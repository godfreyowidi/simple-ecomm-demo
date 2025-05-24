package repo

import (
	"context"
	"fmt"

	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepo struct {
	DB *pgxpool.Pool
}

func NewCustomerRepo(db *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{DB: db}
}

// inserts a new customer
func (r *CustomerRepo) CreateCustomer(ctx context.Context, name string, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.DB.QueryRow(ctx,
		`INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at`,
		name, email,
	).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return &customer, nil
}

// fetches a customer by ID
func (r *CustomerRepo) GetCustomer(ctx context.Context, id int) (*models.Customer, error) {
	var customer models.Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM customers WHERE id = $1`,
		id,
	).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return &customer, nil
}

// fetches a customer by email
func (r *CustomerRepo) GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM customers WHERE email = $1`,
		email,
	).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get customer by email: %w", err)
	}
	return &customer, nil
}

// returns all customers
func (r *CustomerRepo) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, name, email, created_at FROM customers`)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
