package user

import (
	"context"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	IUserRepository interface {
		Register(ctx context.Context, tx *gorm.DB, user User) error
		GetUserByID(ctx context.Context, tx *gorm.DB, userID string) (User, bool, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (User, bool, error)
		GetAllUser(ctx context.Context, tx *gorm.DB, search string) ([]User, error)
		GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req UserPaginationRequest) (UserPaginationRepositoryResponse, error)
		CreateUser(ctx context.Context, tx *gorm.DB, user User) error
		UpdateUser(ctx context.Context, tx *gorm.DB, user User) error
		DeleteUserByID(ctx context.Context, tx *gorm.DB, userID string) error
	}

	UserRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
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


func (ur *UserRepository) Register(ctx context.Context, tx *gorm.DB, user User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Create(&user).Error
}

func (ur *UserRepository) GetUserByID(ctx context.Context, tx *gorm.DB, userID string) (User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var user User
	if err := tx.WithContext(ctx).Where("id = ?", userID).Take(&user).Error; err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var user User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

func (ur *UserRepository) GetAllUser(ctx context.Context, tx *gorm.DB, search string) ([]User, error) {
	if tx == nil {
		tx = ur.db
	}

	var users []User

	query := tx.WithContext(ctx).Model(&User{})

	if search != "" {
		searchValue := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchValue, searchValue)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil

}

func (ur *UserRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req UserPaginationRequest) (UserPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = ur.db
	}

	var users []User
	var err error
	var count int64

	if req.PaginationRequest.PerPage == 0 {
		req.PaginationRequest.PerPage = 10
	}

	if req.PaginationRequest.Page == 0 {
		req.PaginationRequest.Page = 1
	}

	query := tx.WithContext(ctx).Model(&User{})

	if req.PaginationRequest.Search != "" {
		searchValue := "%" + strings.ToLower(req.PaginationRequest.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchValue, searchValue)
	}

	if req.UserID != "" {
		query = query.Where("id = ?", req.UserID)
	}

	if err := query.Count(&count).Error; err != nil {
		return UserPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.PaginationRequest.Page, req.PaginationRequest.PerPage)).Find(&users).Error; err != nil {
		return UserPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PaginationRequest.PerPage)))

	return UserPaginationRepositoryResponse{
		Users: users,
		PaginationResponse: PaginationResponse{
			Page:    req.PaginationRequest.Page,
			PerPage: req.PaginationRequest.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

func (ur *UserRepository) CreateUser(ctx context.Context, tx *gorm.DB, user User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Create(&user).Error
}

func (ur *UserRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Where("id = ?", user.ID).Updates(&user).Error
}

func (ur *UserRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, userID string) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Where("id = ?", userID).Delete(&User{}).Error
}
