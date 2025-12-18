package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/mferdian/Go-GraphQL/config/jwt"
	"github.com/mferdian/Go-GraphQL/constants"
	"github.com/mferdian/Go-GraphQL/logging"
)

type (
	IProductService interface {
		CreateProduct(ctx context.Context, req CreateProductRequest) (ProductResponse, error)
		GetAllProduct(ctx context.Context, search string) ([]ProductResponse, error)
		GetAllProductWithPagination(ctx context.Context, req ProductPaginationRequest) (ProductPaginationResponse, error)
		GetProductByID(ctx context.Context, productID string) (ProductResponse, error)
		UpdateProduct(ctx context.Context, req UpdateProductRequest) (ProductResponse, error)
		DeleteProduct(ctx context.Context, req DeleteProductRequest) (ProductResponse, error)
	}

	ProductService struct {
		productRepo IProductRepository
		jwtService  jwt.InterfaceJWTService
	}
)

func NewProductService(productRepo IProductRepository, jwtService jwt.InterfaceJWTService) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		jwtService:  jwtService,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest) (ProductResponse, error) {
	if len(req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_PRODUCT + ": name too short")
		return ProductResponse{}, constants.ErrInvalidName
	}

	if len(req.Description) < 8 {
		logging.Log.Warn(constants.MESSAGE_FAILED_CREATE_PRODUCT + ": description too short")
		return ProductResponse{}, constants.ErrInvalidDescription
	}

	if req.Price <= 0 {
		return ProductResponse{}, constants.ErrInvalidPrice
	}

	product := Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Merk:        req.Merk,
		Material:    req.Material,
		Price:       req.Price,
	}

	err := ps.productRepo.CreateProduct(ctx, nil, product)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_CREATE_PRODUCT)
		return ProductResponse{}, constants.ErrCretaeProduct
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_CREATE_PRODUCT+": %s", product.Name)

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Merk:        product.Merk,
		Material:    product.Material,
		Price:       product.Price,
	}, nil
}

func (ps *ProductService) GetAllProduct(ctx context.Context, search string) ([]ProductResponse, error) {
	users, err := ps.productRepo.GetAllProduct(ctx, nil, search)

	if err != nil {
		return nil, constants.ErrGetAllProduct
	}

	var datas []ProductResponse
	for _, products := range users {
		data := ProductResponse{
			ID:          products.ID,
			Name:        products.Name,
			Description: products.Description,
			Merk:        products.Merk,
			Material:    products.Material,
			Price:       products.Price,
		}

		datas = append(datas, data)
	}
	return datas, nil
}

func (ps *ProductService) GetAllProductWithPagination(ctx context.Context, req ProductPaginationRequest) (ProductPaginationResponse, error) {
	dataWithPaginate, err := ps.productRepo.GetAllProductWithPagination(ctx, nil, req)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_GET_ALL_PRODUCTS)
		return ProductPaginationResponse{}, constants.ErrGetAllProduct
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_ALL_PRODUCT+": page %d", req.Page)

	var datas []ProductResponse
	for _, product := range dataWithPaginate.Products {
		datas = append(datas, ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Merk:        product.Merk,
			Material:    product.Material,
			Price:       product.Price,
		})
	}

	return ProductPaginationResponse{
		Data: datas,
		PaginationResponse: PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (ps *ProductService) GetProductByID(ctx context.Context, productID string) (ProductResponse, error) {
	if _, err := uuid.Parse(productID); err != nil {
		logging.Log.Warn(constants.MESSAGE_FAILED_GET_DETAIL_USER + ": invalid UUID")
		return ProductResponse{}, constants.ErrInvalidUUID
	}

	product, _, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", productID).Error(constants.MESSAGE_FAILED_GET_DETAIL_USER)
		return ProductResponse{}, constants.ErrGetUserByID
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_GET_DETAIL_USER+": %s", productID)

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Merk:        product.Merk,
		Material:    product.Material,
		Price:       product.Price,
	}, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, req UpdateProductRequest) (ProductResponse, error) {
	product, _, err := ps.productRepo.GetProductByID(ctx, nil, req.ID)
	if err != nil {
		logging.Log.WithError(err).WithField("id", req.ID).Error(constants.MESSAGE_FAILED_UPDATE_PRODUCT)
		return ProductResponse{}, constants.ErrGetProductByID
	}

	if req.Name != nil && len(*req.Name) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_PRODUCT + ": invalid name")
		return ProductResponse{}, constants.ErrInvalidName
	} else if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Description != nil && len(*req.Description) < 5 {
		logging.Log.Warn(constants.MESSAGE_FAILED_UPDATE_PRODUCT + ": invalid name")
		return ProductResponse{}, constants.ErrInvalidDescription
	} else if req.Description != nil {
		product.Description = *req.Description
	}

	if req.Material != nil {
		product.Material = *req.Material
	}

	if req.Merk != nil {
		product.Merk = *req.Merk
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	err = ps.productRepo.UpdateProduct(ctx, nil, product)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_UPDATE_PRODUCT)
		return ProductResponse{}, constants.ErrUpdateProduct
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_UPDATE_PRODUCT+": %s", product.ID)

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Material:    product.Material,
		Merk:        product.Merk,
		Price:       product.Price,
	}, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, req DeleteProductRequest) (ProductResponse, error) {
	product, _, err := ps.productRepo.GetProductByID(ctx, nil, req.ProductID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_PRODUCT)
		return ProductResponse{}, constants.ErrGetProductByID
	}

	err = ps.productRepo.DeleteProduct(ctx, nil, req.ProductID)
	if err != nil {
		logging.Log.WithError(err).Error(constants.MESSAGE_FAILED_DELETE_PRODUCT)
		return ProductResponse{}, constants.ErrDeleteProduct
	}

	logging.Log.Infof(constants.MESSAGE_SUCCESS_DELETE_USER+": %s", req.ProductID)

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Material:    product.Material,
		Merk:        product.Merk,
		Price:       product.Price,
	}, nil
}
