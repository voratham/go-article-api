package migrations

import (
	"article-api/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func m1614447640AddCategoryIdToArticles() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1614447640",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Article{})

			var articles []models.Article
			tx.Find(&articles)

			// just bussiness example force categoryID equal id 2
			for _, article := range articles {
				article.CategoryID = 2
				tx.Save(&article)
			}

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&models.Article{}, "category_id")

		},
	}

}
