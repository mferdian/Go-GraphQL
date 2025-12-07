package user

import "github.com/google/uuid"

type (
	UserResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone_number"`
		Address     string    `json:"address"`
	}

	RegisterUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegisterUserResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	LoginUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	CreateUserRequest struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
	}

	UpdateUserRequest struct {
		ID          string  `json:"-"`
		Name        *string `json:"name,omitempty"`
		Email       *string `json:"email,omitempty"`
		Password    *string `json:"password,omitempty"`
		PhoneNumber *string `json:"phone_number,omitempty"`
		Address     *string `json:"address,omitempty"`
	}

	DeleteUserRequest struct {
		UserID string `json:"-"`
	}

	UserPaginationRequest struct {
		PaginationRequest
		UserID string `form:"id"`
	}

	UserPaginationResponse struct {
		PaginationResponse
		Data []UserResponse `json:"data"`
	}

	UserPaginationRepositoryResponse struct {
		PaginationResponse
		Users []User
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

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (pr *PaginationResponse) GetLimit() int {
	return pr.PerPage
}

func (pr *PaginationResponse) GetPage() int {
	return pr.Page
}
