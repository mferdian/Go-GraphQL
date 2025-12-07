package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/mferdian/Go-GraphQL/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `json:"name"`
	Email       string    `gorm:"unique; not null" json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Role        string    `json:"role"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return nil
}
