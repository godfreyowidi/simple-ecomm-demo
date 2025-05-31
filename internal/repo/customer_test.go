package repo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Fatal("TEST_DATABASE_URL environment variable is not set")
	}

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}

	return dbpool
}

func TestCreateCustomer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repo.NewCustomerRepo(db)

	customer := &models.Customer{
		AuthID:    "auth0|unit-test-" + time.Now().Format("20060102150405.000000000"),
		FirstName: "Test",
		LastName:  "User",
		Email:     "unit_test_" + time.Now().Format("20060102150405") + "@example.com",
		Phone:     "+1234567890",
	}

	created, err := repo.CreateCustomer(context.Background(), customer)
	if err != nil {
		t.Fatalf("CreateCustomer failed: %v", err)
	}

	if created.ID == 0 {
		t.Error("Expected created customer to have a non-zero ID")
	}
	if created.FirstName != customer.FirstName || created.Email != customer.Email {
		t.Errorf("Expected %+v, got %+v", customer, created)
	}
	if created.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}
