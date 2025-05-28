package pkg

import (
	"context"

	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/models"
)

type RegisterInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

type RegisterHandler struct {
	CustomerRepo *repo.CustomerRepo
}

func (h *RegisterHandler) Register(ctx context.Context, input RegisterInput) (*models.Customer, error) {
	token, err := GetManagementToken(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := CreateUser(ctx, token, CreateUserPayload{
		Email:      input.Email,
		Password:   input.Password,
		Connection: "Username-Password-Authentication",
		GivenName:  input.FirstName,
		FamilyName: input.LastName,
	})
	if err != nil {
		return nil, err
	}

	customer := &models.Customer{
		AuthID:    userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
	}

	// Use the actual repo method
	_, err = h.CustomerRepo.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
