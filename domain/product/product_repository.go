package product

import (
	"context"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	IProductRepository interface {
		CreateProduct(ctx context.Context, tx *gorm.DB, product Product) error
		GetProductByID(ctx context.Context, tx *gorm.DB, productID string) (Product, bool, error)
		GetAllProduct(ctx context.Context, tx *gorm.DB, search string) ([]Product, error)
		GetAllProductWithPagination(ctx context.Context, tx *gorm.DB, req ProductPaginationRequest) (ProductPaginationRepositoryResponse, error)
		UpdateProduct(ctx context.Context, tx *gorm.DB, product Product) error
		DeleteProduct(ctx context.Context, tx *gorm.DB, productID string) error
	}

	ProductRepository struct {
		db *gorm.DB
	}
)

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Pagination
func Paginate(page, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, tx *gorm.DB, product Product) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&product).Error
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, tx *gorm.DB, productID string) (Product, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var product Product
	if err := tx.WithContext(ctx).Where("id = ?", productID).Take(&product).Error; err != nil {
		return Product{}, false, err
	}

	return product, true, nil
}

func (pr *ProductRepository) GetAllProduct(ctx context.Context, tx *gorm.DB, search string) ([]Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var users []Product

	query := tx.WithContext(ctx).Model(&Product{})

	if search != "" {
		searchValue := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(merk) LIKE ? OR LOWER(material) LIKE ?",
			searchValue, searchValue, searchValue)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (pr *ProductRepository) GetAllProductWithPagination(ctx context.Context, tx *gorm.DB, req ProductPaginationRequest) (ProductPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var product []Product
	var err error
	var count int64

	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}

	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).Model(&Product{})

	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchValue, searchValue)
	}

	if req.UserID != "" {
		query = query.Where("id = ?", req.UserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return ProductPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&product).Error; err != nil {
		return ProductPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return ProductPaginationRepositoryResponse{
		Products: product,
		PaginationResponse: PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, tx *gorm.DB, product Product) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", product.ID).Updates(&product).Error
}
func (pr *ProductRepository) DeleteProduct(ctx context.Context, tx *gorm.DB, productID string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", productID).Delete(&Product{}).Error
}
