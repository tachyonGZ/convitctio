package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:50"`
	Password string `json:"-"`
}

func GetUserByID() User {
	var u User
	return u
}

func NewUser() User {
	return User{}
}
