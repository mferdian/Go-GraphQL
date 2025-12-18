package migrations

import (
	"github.com/mferdian/Go-GraphQL/domain/product"
	"github.com/mferdian/Go-GraphQL/domain/user"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	err := SeedFromJSON[user.User](db, "./migrations/json/users.json", user.User{}, "Email")
	if err != nil {
		return err
	}
	err = SeedFromJSON[product.Product](db, "./migrations/json/products.json", product.Product{}, "Merk")
	if err != nil {
		return err
	}

	return nil
}