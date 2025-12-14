package product

import "github.com/google/uuid"

type (
	// GraphQL
	CreateProductInput struct {
		Name        string
		Description string
		Merk        string
		Material    string
		Price       float64
	}

	ProductResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Merk        string    `json:"merk"`
		Material    string    `json:"material"`
		Price       float32   `json:"price"`
	}

	CreateProductRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Merk        string  `json:"merk"`
		Material    string  `json:"material"`
		Price       float32 `json:"price"`
	}

	UpdateProductRequest struct {
		ID          string   `json:"-"`
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Merk        *string  `json:"merk"`
		Material    *string  `json:"material"`
		Price       *float32 `json:"price"`
	}

	DeleteProductRequest struct {
		ProductID string `json:"-"`
	}

	ProductPaginationRequest struct {
		PaginationRequest
		UserID string `form:"id"`
	}

	ProductPaginationResponse struct {
		PaginationResponse
		Data []ProductResponse `json:"data"`
	}

	ProductPaginationRepositoryResponse struct {
		PaginationResponse
		Products []Product
	}

	PaginationRequest struct {
		Search  string `form:"search"`
		Page    int    `form:"page"`
		PerPage int    `form:"per_page"`
	}

	PaginationResponse struct {
		Page    int   `json:"page"`
		PerPage int   `json:"per_page"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)
