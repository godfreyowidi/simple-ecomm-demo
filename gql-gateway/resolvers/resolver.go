package resolvers

import (
	"github.com/godfreyowidi/simple-ecomm-demo/internal/repo"
	"github.com/godfreyowidi/simple-ecomm-demo/pkg"
)

type Resolver struct {
	ProductRepo     *repo.ProductRepo
	CustomerRepo    *repo.CustomerRepo
	OrderRepo       *repo.OrderRepo
	OrderItemRepo   *repo.OrderItemRepo
	CategoryRepo    *repo.CategoryRepo
	RegisterHandler *pkg.RegisterHandler
	CatalogRepo     *repo.CatalogRepo
}
