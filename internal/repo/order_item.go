package repo

import (
	"context"
	"fmt"

	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemRepo struct {
	DB *pgxpool.Pool
}

func NewOrderItemRepo(db *pgxpool.Pool) *OrderItemRepo {
	return &OrderItemRepo{DB: db}
}

// inserts an item into an order
func (r *OrderItemRepo) CreateOrderItem(ctx context.Context, orderID, productID, quantity int, price float64) (*models.OrderItem, error) {
	var item models.OrderItem
	err := r.DB.QueryRow(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price)
		 VALUES ($1, $2, $3, $4) RETURNING id, order_id, product_id, quantity, price`,
		orderID, productID, quantity, price,
	).Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
	if err != nil {
		return nil, fmt.Errorf("create order item: %w", err)
	}
	return &item, nil
}

// get all items for a specific order
func (r *OrderItemRepo) GetItemsByOrder(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, order_id, product_id, quantity, price
		 FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("get items by order: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
