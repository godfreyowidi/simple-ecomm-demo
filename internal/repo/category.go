package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID       int
	Name     string
	ParentID *int
}

type CategoryRepo struct {
	DB *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{DB: db}
}

// inserts a new category
func (r *CategoryRepo) CreateCategory(ctx context.Context, name string, parentID *int) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING id`,
		name, parentID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create category: %w", err)
	}
	return id, nil
}

// fetches a category by ID
func (r *CategoryRepo) GetCategory(ctx context.Context, id int) (*Category, error) {
	var c Category
	err := r.DB.QueryRow(ctx,
		`SELECT id, name, parent_id FROM categories WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.ParentID)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	return &c, nil
}

// returns all categories
func (r *CategoryRepo) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, name, parent_id FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
