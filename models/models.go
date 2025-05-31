package models

import "time"

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id,omitempty"`
}

type Customer struct {
	ID        int       `json:"id"`
	AuthID    string    `json:"auth_id"` // Auth0 user ID (e.g. "auth0|abc123")
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemInput struct {
	ProductID int
	Quantity  int
	Price     float64
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	OrderDate  time.Time `json:"order_date"`
	Status     string    `json:"status"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
	CategoryID  *int    `json:"category_id,omitempty"`
}

type ProductCatalog struct {
	TopCategoryName string
	SubCategories   []ProductSubCategory
}

type ProductSubCategory struct {
	Name     string
	Products []Product
}
