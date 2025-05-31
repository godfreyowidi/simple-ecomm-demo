package repo_test

import (
	"context"
	"testing"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
)

func TestCreateCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repo.NewCategoryRepo(db)

	// top-level category
	cat, err := repo.CreateCategory(context.Background(), "Test Category", nil)
	if err != nil {
		t.Fatalf("CreateCategory failed: %v", err)
	}

	if cat.ID == 0 {
		t.Error("Expected created category to have a non-zero ID")
	}
	if cat.Name != "Test Category" {
		t.Errorf("Expected name to be 'Test Category', got '%s'", cat.Name)
	}
	if cat.ParentID != nil {
		t.Errorf("Expected ParentID to be nil, got %v", *cat.ParentID)
	}

	// sub-category
	subCat, err := repo.CreateCategory(context.Background(), "Sub Category", &cat.ID)
	if err != nil {
		t.Fatalf("Create sub-category failed: %v", err)
	}
	if subCat.ParentID == nil || *subCat.ParentID != cat.ID {
		t.Errorf("Expected ParentID to be %d, got %v", cat.ID, subCat.ParentID)
	}
}
