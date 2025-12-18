package resolver

import "github.com/mferdian/Go-GraphQL/domain/product"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ProductService product.IProductService
}