package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/godfreyowidi/simple-ecomm-demo/gql-gateway/graph"
	"github.com/godfreyowidi/simple-ecomm-demo/gql-gateway/resolvers"
	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/migrations"
	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newGraphQLServer(pool *pgxpool.Pool) *handler.Server {
	res := &resolvers.Resolver{
		ProductRepo:   repo.NewProductRepo(pool),
		CustomerRepo:  repo.NewCustomerRepo(pool),
		OrderRepo:     repo.NewOrderRepo(pool),
		OrderItemRepo: repo.NewOrderItemRepo(pool),
		CategoryRepo:  repo.NewCategoryRepo(pool),
	}
	return handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: res}))
}

func setupPostgres(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Fatal("TEST_DATABASE_URL environment variable is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if err := migrations.UpWithDSN(ctx, connStr); err != nil {
		pool.Close()
		t.Fatalf("migrations failed: %v", err)
	}

	teardown := func() {
		pool.Close()
	}

	return pool, teardown
}

func TestCreateOrderMutation(t *testing.T) {
	pool, teardown := setupPostgres(t)
	defer teardown()

	ctx := context.Background()

	// Create category
	cat, err := repo.NewCategoryRepo(pool).CreateCategory(ctx, "Elec", nil)
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	// Create product
	prod, err := repo.NewProductRepo(pool).CreateProduct(ctx, "Phone", nil, 699.99, &cat.ID)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	// Create customer with full details
	customer := &models.Customer{
		AuthID:    "auth0|testid",
		FirstName: "Bob",
		LastName:  "Smith",
		Email:     uniqueEmail(),
		Phone:     "+123456789",
	}
	createdCustomer, err := repo.NewCustomerRepo(pool).CreateCustomer(ctx, customer)
	if err != nil {
		t.Fatalf("create customer: %v", err)
	}

	// Start GraphQL server
	srv := newGraphQLServer(pool)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	query := `
	mutation ($cid: ID!, $pid: ID!) {
		createOrder(input: {
			customerID: $cid,
			items: [{ productID: $pid, quantity: 2, price: 699.99 }]
		}) {
			id
			status
			items {
				quantity
			}
		}
	}`

	reqBody := map[string]any{
		"query": query,
		"variables": map[string]any{
			"cid": fmt.Sprintf("%d", createdCustomer.ID), // send as string
			"pid": fmt.Sprintf("%d", prod.ID),
		},
	}
	b, _ := json.Marshal(reqBody)
	resp, err := ts.Client().Post(ts.URL, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	var out struct {
		Data struct {
			CreateOrder struct {
				ID     string
				Status string
				Items  []struct {
					Quantity int
				}
			}
		}
		Errors json.RawMessage
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(out.Errors) != 0 {
		t.Fatalf("got errors: %s", string(out.Errors))
	}
	if out.Data.CreateOrder.Status != "pending" || out.Data.CreateOrder.Items[0].Quantity != 2 {
		t.Errorf("unexpected result: %+v", out.Data.CreateOrder)
	}
}

func TestUploadProductsWithCategoriesForGQL(t *testing.T) {
	pool, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	parentCat, err := repo.NewCategoryRepo(pool).CreateCategory(ctx, "Electronics", nil)
	if err != nil {
		t.Fatalf("create parent category: %v", err)
	}
	subCat, err := repo.NewCategoryRepo(pool).CreateCategory(ctx, "Smartphones", &parentCat.ID)
	if err != nil {
		t.Fatalf("create sub-category: %v", err)
	}

	prod, err := repo.NewProductRepo(pool).CreateProduct(ctx, "iPhone 13", nil, 999.99, &subCat.ID)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	if prod.Name != "iPhone 13" || prod.CategoryID == nil || *prod.CategoryID != subCat.ID {
		t.Errorf("unexpected product: %+v", prod)
	}
}

func TestAverageProductPriceByCategoryForGQL(t *testing.T) {
	pool, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	cat, err := repo.NewCategoryRepo(pool).CreateCategory(ctx, "Laptops", nil)
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	_, err = repo.NewProductRepo(pool).CreateProduct(ctx, "Laptop A", nil, 1000.00, &cat.ID)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}
	_, err = repo.NewProductRepo(pool).CreateProduct(ctx, "Laptop B", nil, 500.00, &cat.ID)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	avg, err := repo.NewProductRepo(pool).GetAveragePriceByCategory(ctx, cat.ID)
	if err != nil {
		t.Fatalf("get avg price: %v", err)
	}

	expected := (1000.00 + 500.00) / 2
	if avg != expected {
		t.Errorf("expected %v, got %v", expected, avg)
	}
}

func TestMakeOrdersForGQL(t *testing.T) {
	pool, teardown := setupPostgres(t)
	defer teardown()
	ctx := context.Background()

	cat, err := repo.NewCategoryRepo(pool).CreateCategory(ctx, "Accessories", nil)
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	p1, err := repo.NewProductRepo(pool).CreateProduct(ctx, "Mouse", nil, 25.00, &cat.ID)
	if err != nil {
		t.Fatalf("create product 1: %v", err)
	}
	p2, err := repo.NewProductRepo(pool).CreateProduct(ctx, "Keyboard", nil, 45.00, &cat.ID)
	if err != nil {
		t.Fatalf("create product 2: %v", err)
	}

	cust, err := repo.NewCustomerRepo(pool).CreateCustomer(ctx, &models.Customer{
		AuthID:    "auth0|alice123",
		FirstName: "Alice",
		LastName:  "Test",
		Email:     uniqueEmail(),
		Phone:     "+1234567890",
	})
	if err != nil {
		t.Fatalf("create customer: %v", err)
	}

	items := []models.OrderItemInput{
		{ProductID: p1.ID, Quantity: 2, Price: 25.00},
		{ProductID: p2.ID, Quantity: 1, Price: 45.00},
	}

	order, err := repo.NewOrderRepo(pool).CreateOrder(ctx, cust.ID, items)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if order.ID == 0 || order.CustomerID != cust.ID {
		t.Errorf("unexpected order: %+v", order)
	}
}

func uniqueEmail() string {
	return fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
}
