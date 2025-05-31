package repo

import (
	"context"
	"fmt"

	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{DB: db}
}

// inserts a new product
func (r *ProductRepo) CreateProduct(ctx context.Context, name string, description *string, price float64, categoryID *int) (*models.Product, error) {
	var p models.Product
	err := r.DB.QueryRow(ctx,
		`INSERT INTO products (name, description, price, category_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, description, price, category_id`,
		name, description, price, categoryID,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return &p, nil
}

// get a product by ID
func (r *ProductRepo) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	var p models.Product
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
func (r *ProductRepo) ListProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, name, description, price, category_id FROM products`)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// returns the average price of products in a category
func (r *ProductRepo) GetAveragePriceByCategory(ctx context.Context, categoryID int) (float64, error) {
	var avg float64
	err := r.DB.QueryRow(ctx,
		`SELECT AVG(price) FROM products WHERE category_id = $1`,
		categoryID,
	).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("get average price: %w", err)
	}
	return avg, nil
}
