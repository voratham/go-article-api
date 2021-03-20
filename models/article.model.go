package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title      string `gorm:"unique;not null"`
	Excerpt    string `gorm:"not null"`
	Body       string `gorm:"not null"`
	Image      string `gorm:"not null"`
	CategoryID int
	// gorm library checking foreignKey has on category table
	Category Category `gorm:"foreignKey:CategoryID"`
	UserID   uint
	User     User
}
