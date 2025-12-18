package migrations

import (
	"github.com/mferdian/Go-GraphQL/domain/product"
	"github.com/mferdian/Go-GraphQL/domain/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&user.User{},
		&product.Product{},
	); err != nil {
		return err
	}

	return nil
}
