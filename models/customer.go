package models

import "time"

type Customer struct {
	ID        int       `json:"id"`
	AuthID    string    `json:"auth_id"` // Auth0 user ID (e.g. "auth0|abc123")
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
