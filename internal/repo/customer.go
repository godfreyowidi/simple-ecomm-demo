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

// new customer
func (r *CustomerRepo) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	err := r.DB.QueryRow(ctx,
		`INSERT INTO customers (auth_id, first_name, last_name, email, phone)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, auth_id, first_name, last_name, email, phone, created_at`,
		customer.AuthID, customer.FirstName, customer.LastName, customer.Email, customer.Phone,
	).Scan(
		&customer.ID,
		&customer.AuthID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.Phone,
		&customer.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return customer, nil
}

// get a customer by ID
func (r *CustomerRepo) GetCustomerById(ctx context.Context, id int) (*models.Customer, error) {
	var c models.Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, auth_id, first_name, last_name, email, phone, created_at
		 FROM customers WHERE id = $1`, id,
	).Scan(
		&c.ID, &c.AuthID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return &c, nil
}

// get a customer by email
func (r *CustomerRepo) GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var c models.Customer
	err := r.DB.QueryRow(ctx,
		`SELECT id, auth_id, first_name, last_name, email, phone, created_at
		 FROM customers WHERE email = $1`, email,
	).Scan(
		&c.ID, &c.AuthID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get customer by email: %w", err)
	}
	return &c, nil
}

// get all customers
func (r *CustomerRepo) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, auth_id, first_name, last_name, email, phone, created_at FROM customers`)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.AuthID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

// get customer by email/phone
func (r *CustomerRepo) FindByEmailOrPhone(ctx context.Context, identifier string) (*models.Customer, error) {
	query := `
		SELECT id, email, phone FROM customers 
		WHERE email = $1 OR phone = $1 LIMIT 1;
	`
	var c models.Customer
	err := r.DB.QueryRow(ctx, query, identifier).Scan(&c.ID, &c.Email, &c.Phone)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
