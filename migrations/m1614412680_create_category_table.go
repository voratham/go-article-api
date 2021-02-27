package migrations

import (
	"article-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1614412680CreateCategoryTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1614412680",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("categoires")

		},
	}

}
