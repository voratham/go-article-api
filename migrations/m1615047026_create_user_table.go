package migrations

import (
	"article-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1615047026CreateUserTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1615047026",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")

		},
	}

}
