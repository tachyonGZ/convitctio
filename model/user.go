package model

import (
	"conviction/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:50"`
	Password string `json:"-"`
}

func GetUserByID(ID any) (User, error) {
	var user User
	result := db.GetDB().First(&user, ID)
	return user, result.Error
}

func GetUserByUsername(username string) (User, error) {
	var user User
	result := db.GetDB().Where("username = ?", username).First(&user)
	return user, result.Error
}

func NewUser() User {
	return User{}
}

func (user *User) CheckPassword(password string) bool {
	return password == user.Password
}
