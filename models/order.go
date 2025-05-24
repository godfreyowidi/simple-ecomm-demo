package models

import "time"

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	OrderDate  time.Time `json:"order_date"`
	Status     string    `json:"status"`
}
