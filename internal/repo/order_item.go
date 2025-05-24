package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
}

type OrderItemRepo struct {
	DB *pgxpool.Pool
}

func NewOrderItemRepo(db *pgxpool.Pool) *OrderItemRepo {
	return &OrderItemRepo{DB: db}
}

// inserts an item into an order
func (r *OrderItemRepo) CreateOrderItem(ctx context.Context, orderID, productID, quantity int, price float64) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		orderID, productID, quantity, price,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create order item: %w", err)
	}
	return id, nil
}

// retrieves all items for a specific order
func (r *OrderItemRepo) GetItemsByOrder(ctx context.Context, orderID int) ([]OrderItem, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, order_id, product_id, quantity, price
		 FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("get items by order: %w", err)
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
