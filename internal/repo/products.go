package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID          int
	Name        string
	Description *string
	Price       float64
	CategoryID  *int
}

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{DB: db}
}

// inserts a new product
func (r *ProductRepo) CreateProduct(ctx context.Context, name string, description *string, price float64, categoryID *int) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO products (name, description, price, category_id)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		name, description, price, categoryID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create product: %w", err)
	}
	return id, nil
}

// fetches a product by ID
func (r *ProductRepo) GetProduct(ctx context.Context, id int) (*Product, error) {
	var p Product
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, description, price, category_id FROM products WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return &p, nil
}

// returns all products
func (r *ProductRepo) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, name, description, price, category_id FROM products`)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
