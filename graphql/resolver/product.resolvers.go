package resolver

import (
	"context"

	"github.com/mferdian/Go-GraphQL/domain/product"
	"github.com/mferdian/Go-GraphQL/graphql/generated"
	"github.com/mferdian/Go-GraphQL/graphql/model"
)

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, search *string) ([]*model.Product, error) {
	var keyword string
	if search != nil {
		keyword = *search
	}

	products, err := r.ProductService.GetAllProduct(ctx, keyword)
	if err != nil {
		return nil, err
	}

	var result []*model.Product
	for _, p := range products {
		result = append(result, &model.Product{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Merk:        &p.Merk,
			Material:    &p.Material,
			Price:       float64(p.Price),
		})
	}

	return result, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	p, err := r.ProductService.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Merk:        &p.Merk,
		Material:    &p.Material,
		Price:       float64(p.Price),
	}, nil
}

// ProductsWithPagination is the resolver for the productsWithPagination field.
func (r *queryResolver) ProductsWithPagination(ctx context.Context, page int, perPage int, search *string) (*model.ProductPagination, error) {
	req := product.ProductPaginationRequest{
		PaginationRequest: product.PaginationRequest{
			Page:    page,
			PerPage: perPage,
		},
	}

	if search != nil {
		req.Search = *search
	}

	data, err := r.ProductService.GetAllProductWithPagination(ctx, req)
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	for _, p := range data.Data {
		products = append(products, &model.Product{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Merk:        &p.Merk,
			Material:    &p.Material,
			Price:       float64(p.Price),
		})
	}

	return &model.ProductPagination{
		Data: products,
		Pagination: &model.Pagination{
			Page:    data.Page,
			PerPage: data.PerPage,
			MaxPage: int(data.MaxPage),
			Count:   int(data.Count),
		},
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
