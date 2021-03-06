package migrations

import (
	"article-api/config"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			// add script migration
			m1614330748CreateArticlesTable(),
			m1614412680CreateCategoryTable(),
			m1614447640AddCategoryIdToArticles(),
			m1616223041AddUserIdToArticles(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatal("Could not migrate: %v", err)
	}

}
