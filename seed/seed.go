package seed

import (
	"article-api/config"
	"article-api/migrations"
	"article-api/models"
	"math/rand"
	"strconv"

	"log"

	"github.com/bxcodec/faker/v3"
)

func Load() {
	db := config.GetDB()

	db.Migrator().DropTable("articles", "categories", "users", "migrations")
	migrations.Migrate()

	// Add categories
	log.Println("Creating categories...")

	numOfCategoires := 20
	categories := make([]models.Category, 0, numOfCategoires)

	for i := 1; i <= numOfCategoires; i++ {

		category := models.Category{
			Name: faker.Word() + "#" + strconv.Itoa(i),
			Desc: faker.Paragraph(),
		}

		db.Create(&category)
		categories = append(categories, category)

	}

	// Add articles
	log.Println("Creating articles...")

	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)

	for i := 1; i <= numOfArticles; i++ {

		article := models.Article{
			Title:      "#" + strconv.Itoa(i) + faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unspash.com/random/300*200?" + strconv.Itoa(i),
			CategoryID: rand.Intn(numOfCategoires) + 1,
		}

		db.Create(&article)

		articles = append(articles, article)

	}

}
