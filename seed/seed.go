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

	db.Migrator().DropTable("users", "articles", "categories", "migrations")
	migrations.Migrate()

	//Add Admin
	log.Println("Creating admin...")
	admin := models.User{
		Email:    "admin@admin.com",
		Name:     "admin",
		Password: "12345678",
		Avatar:   "https://i.pravatar.cc/100",
		Role:     "Admin",
	}

	admin.Password = admin.GenerateEncryptedPassword()
	db.Create(&admin)

	// Add users
	log.Println("Creating users...")

	numOfUsers := 20
	users := make([]models.User, 0, numOfUsers)
	usersRoles := [2]string{"Editor", "Member"}

	for i := 1; i <= numOfUsers; i++ {
		user := models.User{
			Email:    faker.Email(),
			Name:     faker.Name(),
			Password: "12345678",
			Avatar:   "https://i.pravatar.cc/100?" + strconv.Itoa(i),
			Role:     usersRoles[rand.Intn(2)],
		}
		user.Password = user.GenerateEncryptedPassword()
		db.Create(&user)
		users = append(users, user)
	}

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
			Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			CategoryID: rand.Intn(numOfCategoires) + 1,
			UserID:     uint(rand.Intn(numOfUsers) + 1),
		}

		db.Create(&article)

		articles = append(articles, article)

	}

}
