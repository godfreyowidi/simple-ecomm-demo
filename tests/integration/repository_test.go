package integration

import (
	"context"
	"os"
	"testing"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/migrations"
	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestRepos struct {
	Pool          *pgxpool.Pool
	ProductRepo   *repo.ProductRepo
	CustomerRepo  *repo.CustomerRepo
	OrderRepo     *repo.OrderRepo
	OrderItemRepo *repo.OrderItemRepo
	CategoryRepo  *repo.CategoryRepo
}

func setupPostgresAndRepos(t *testing.T) (*TestRepos, func()) {
	t.Helper()
	ctx := context.Background()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Fatal("TEST_DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("connect to DB: %v", err)
	}

	if err := migrations.UpWithDSN(ctx, connStr); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repos := &TestRepos{
		Pool:          pool,
		ProductRepo:   repo.NewProductRepo(pool),
		CustomerRepo:  repo.NewCustomerRepo(pool),
		OrderRepo:     repo.NewOrderRepo(pool),
		OrderItemRepo: repo.NewOrderItemRepo(pool),
		CategoryRepo:  repo.NewCategoryRepo(pool),
	}

	teardown := func() {
		repos.Pool.Close()
	}

	return repos, teardown
}

func TestUploadProductsWithCategories(t *testing.T) {
	repos, teardown := setupPostgresAndRepos(t)
	defer teardown()
	ctx := context.Background()

	parentCat, err := repos.CategoryRepo.CreateCategory(ctx, "Electronics", nil)
	if err != nil {
		t.Fatalf("create parent category: %v", err)
	}
	subCat, err := repos.CategoryRepo.CreateCategory(ctx, "Smartphones", &parentCat.ID)
	if err != nil {
		t.Fatalf("create sub-category: %v", err)
	}

	prod, err := repos.ProductRepo.CreateProduct(ctx, "iPhone 13", nil, 999.99, &subCat.ID)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	if prod.Name != "iPhone 13" || prod.CategoryID == nil || *prod.CategoryID != subCat.ID {
		t.Errorf("unexpected product: %+v", prod)
	}
}

func TestAverageProductPriceByCategory(t *testing.T) {
	repos, teardown := setupPostgresAndRepos(t)
	defer teardown()
	ctx := context.Background()

	cat, _ := repos.CategoryRepo.CreateCategory(ctx, "Laptops", nil)
	repos.ProductRepo.CreateProduct(ctx, "Laptop A", nil, 1000.00, &cat.ID)
	repos.ProductRepo.CreateProduct(ctx, "Laptop B", nil, 500.00, &cat.ID)

	avg, err := repos.ProductRepo.GetAveragePriceByCategory(ctx, cat.ID)
	if err != nil {
		t.Fatalf("get avg price: %v", err)
	}

	expected := (1000.00 + 500.00) / 2
	if avg != expected {
		t.Errorf("expected %v, got %v", expected, avg)
	}
}

func TestCreateOrderFlow(t *testing.T) {
	repos, teardown := setupPostgresAndRepos(t)
	defer teardown()
	ctx := context.Background()

	cat, _ := repos.CategoryRepo.CreateCategory(ctx, "Accessories", nil)
	p1, _ := repos.ProductRepo.CreateProduct(ctx, "Mouse", nil, 25.00, &cat.ID)
	p2, _ := repos.ProductRepo.CreateProduct(ctx, "Keyboard", nil, 45.00, &cat.ID)

	cust, err := repos.CustomerRepo.CreateCustomer(ctx, "Bob", "bob@example.com")
	if err != nil {
		t.Fatalf("create customer: %v", err)
	}

	items := []models.OrderItemInput{
		{ProductID: p1.ID, Quantity: 2, Price: 25.00},
		{ProductID: p2.ID, Quantity: 1, Price: 45.00},
	}
	order, err := repos.OrderRepo.CreateOrder(ctx, cust.ID, items)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if order.ID == 0 || order.CustomerID != cust.ID {
		t.Errorf("unexpected order: %+v", order)
	}

	fetched, err := repos.OrderRepo.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("fetch order: %v", err)
	}
	if fetched.Status == "" {
		t.Errorf("invalid order status: %+v", fetched)
	}
}
