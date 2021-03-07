package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Avatar   string
	Role     string `gorm:"default:'Member'; not null"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {

	if u.Password == "" {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	tx.Statement.SetColumn("password", string(hash))
	return
}
