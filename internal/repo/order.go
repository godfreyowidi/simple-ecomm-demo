package repo

import (
	"context"
	"fmt"

	"github.com/godfreyowidi/simple-ecomm-demo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	DB *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{DB: db}
}

// insert an order with items
func (r *OrderRepo) CreateOrder(ctx context.Context, customerID int, items []models.OrderItemInput) (*models.Order, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var order models.Order
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (customer_id) VALUES ($1) RETURNING id, customer_id, order_date, status`,
		customerID,
	).Scan(&order.ID, &order.CustomerID, &order.OrderDate, &order.Status)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price)
			 VALUES ($1, $2, $3, $4)`,
			order.ID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &order, nil
}

// get an order by ID
func (r *OrderRepo) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	var o models.Order
	err := r.DB.QueryRow(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	return &o, nil
}

// get all orders
func (r *OrderRepo) ListOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, customer_id, order_date, status FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// get all orders made by a specific customer
func (r *OrderRepo) ListOrdersByCustomer(ctx context.Context, customerID int) ([]models.Order, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE customer_id = $1`,
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders by customer: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// updates the status of a given order
func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	cmdTag, err := r.DB.Exec(ctx,
		`UPDATE orders SET status = $1 WHERE id = $2`,
		status, orderID,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no order found with id %d", orderID)
	}
	return nil
}
