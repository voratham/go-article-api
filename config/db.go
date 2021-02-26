package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {
	log.Println("Init database")
	var err error

	baseURI := os.Getenv("DB_CONNECTION")
	db, err = gorm.Open(postgres.Open(baseURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal(err)

	}
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	log.Println("Close database")
	conn, _ := db.DB()
	conn.Close()

}
