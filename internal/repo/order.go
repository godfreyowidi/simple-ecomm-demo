package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Order struct {
	ID         int
	CustomerID int
	OrderDate  time.Time
	Status     string
}

type OrderRepo struct {
	DB *pgx.Conn
}

func NewOrderRepo(conn *pgx.Conn) *OrderRepo {
	return &OrderRepo{DB: conn}
}

// inserts a new order
func (r *OrderRepo) CreateOrder(ctx context.Context, customerID int, status string) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx,
		`INSERT INTO orders (customer_id, status) VALUES ($1, $2) RETURNING id`,
		customerID, status,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create order: %w", err)
	}
	return id, nil
}

// retrieves an order by ID
func (r *OrderRepo) GetOrder(ctx context.Context, id int) (*Order, error) {
	var o Order
	err := r.DB.QueryRow(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	return &o, nil
}

// lists all orders
func (r *OrderRepo) ListOrders(ctx context.Context) ([]Order, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, customer_id, order_date, status FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.OrderDate, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// returns all orders made by a specific customer
func (r *OrderRepo) ListOrdersByCustomer(ctx context.Context, customerID int) ([]Order, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, customer_id, order_date, status FROM orders WHERE customer_id = $1`,
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders by customer: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
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
