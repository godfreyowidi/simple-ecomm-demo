package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
	CategoryID  *int    `json:"category_id,omitempty"`
}
