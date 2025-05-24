package resolvers

import "github.com/godfreyowidi/simple-ecomm-demo/internal/repo"

type Resolver struct {
	ProductRepo   *repo.ProductRepo
	CustomerRepo  *repo.CustomerRepo
	OrderRepo     *repo.OrderRepo
	OrderItemRepo *repo.OrderItemRepo
}
