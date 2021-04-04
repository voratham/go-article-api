package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {
	log.Println("Init database")
	var err error
	var logLevel logger.LogLevel

	if gin.Mode() == gin.DebugMode {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	baseURI := os.Getenv("DB_CONNECTION")
	db, err = gorm.Open(postgres.Open(baseURI), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
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
