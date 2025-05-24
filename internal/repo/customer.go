package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Customer struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

type CustomerRepo struct {
	DB *pgxpool.Pool
}

func NewCustomerRepo(db *pgxpool.Pool) *CustomerRepo {
	return &CustomerRepo{DB: db}
}

// inserts a new customer
func (r *CustomerRepo) CreateCustomer(ctx context.Context, name string, email string) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id`,
		name, email,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create customer: %w", err)
	}
	return id, nil
}

// fetches a customer by ID
func (r *CustomerRepo) GetCustomer(ctx context.Context, id int) (*Customer, error) {
	var c Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM customers WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return &c, nil
}

// fetches a customer by email
func (r *CustomerRepo) GetCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	var c Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM customers WHERE email = $1`,
		email,
	).Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get customer by email: %w", err)
	}
	return &c, nil
}

// returns all customers
func (r *CustomerRepo) ListCustomers(ctx context.Context) ([]Customer, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, name, email, created_at FROM customers`)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
