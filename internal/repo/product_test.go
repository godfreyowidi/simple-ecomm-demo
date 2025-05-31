package repo_test

import (
	"context"
	"testing"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
)

func TestCreateProduct(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	categoryRepo := repo.NewCategoryRepo(db)
	productRepo := repo.NewProductRepo(db)

	// let create a category first
	category, err := categoryRepo.CreateCategory(ctx, "Test Category", nil)
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	description := "A test product"
	price := 49.99

	product, err := productRepo.CreateProduct(ctx, "Test Product", &description, price, &category.ID)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	if product.ID == 0 {
		t.Error("Expected created product to have a non-zero ID")
	}
	if product.Name != "Test Product" || product.Price != price {
		t.Errorf("Unexpected product data: %+v", product)
	}
	if product.CategoryID == nil || *product.CategoryID != category.ID {
		t.Errorf("Expected product to have category ID %d, got %v", category.ID, product.CategoryID)
	}
}

func TestGetAveragePriceByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	categoryRepo := repo.NewCategoryRepo(db)
	productRepo := repo.NewProductRepo(db)

	// create a test category
	category, err := categoryRepo.CreateCategory(ctx, "Test Category Avg", nil)
	if err != nil {
		t.Fatalf("create category failed: %v", err)
	}

	// multiple products in that category - insert
	_, err = productRepo.CreateProduct(ctx, "Product A", nil, 100.00, &category.ID)
	if err != nil {
		t.Fatalf("create product A failed: %v", err)
	}
	_, err = productRepo.CreateProduct(ctx, "Product B", nil, 200.00, &category.ID)
	if err != nil {
		t.Fatalf("create product B failed: %v", err)
	}

	avg, err := productRepo.GetAveragePriceByCategory(ctx, category.ID)
	if err != nil {
		t.Fatalf("GetAveragePriceByCategory failed: %v", err)
	}

	expected := (100.00 + 200.00) / 2
	if avg != expected {
		t.Errorf("Expected average %v, got %v", expected, avg)
	}
}
