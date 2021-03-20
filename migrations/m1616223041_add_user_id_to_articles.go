package migrations

import (
	"article-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1616223041AddUserIdToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1616223041",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Article{})
			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Article{}, "user_id")

		},
	}
}
