package repo_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/models"
)

func TestCreateOrder(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	customerRepo := repo.NewCustomerRepo(db)
	productRepo := repo.NewProductRepo(db)
	categoryRepo := repo.NewCategoryRepo(db)
	orderRepo := repo.NewOrderRepo(db)

	// create customer
	// create customer
	customer, err := customerRepo.CreateCustomer(ctx, &models.Customer{
		AuthID:    "auth0|order-test-" + RandString(8),
		FirstName: "Order",
		LastName:  "Tester",
		Email:     "order_tester_" + RandString(8) + "@example.com",
		Phone:     "+1111111111",
	})

	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	// create category
	category, err := categoryRepo.CreateCategory(ctx, "Order Test Category", nil)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	// create product
	product, err := productRepo.CreateProduct(ctx, "Test Product", nil, 49.99, &category.ID)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// create order
	order, err := orderRepo.CreateOrder(ctx, customer.ID, []models.OrderItemInput{
		{
			ProductID: product.ID,
			Quantity:  2,
			Price:     product.Price,
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	// Assertions -- checks
	if order.ID == 0 {
		t.Error("expected order to have a non-zero ID")
	}
	if order.CustomerID != customer.ID {
		t.Errorf("expected customer ID %d, got %d", customer.ID, order.CustomerID)
	}
	if order.Status == "" {
		t.Error("expected order to have a status set")
	}
	if order.OrderDate.IsZero() {
		t.Error("expected order date to be set")
	}
}

// Helper

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
