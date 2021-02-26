package migrations

import (
	"article-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1614330748CreateArticlesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1614330748",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("articles")

		},
	}

}
