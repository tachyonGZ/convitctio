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

func (user *User) CheckPassword(password string) bool {
	return password == user.Password
}

func (user *User) Root() (*Directory, error) {
	pRootDir := &Directory{}
	res := db.GetDB().Where("parent_id is NULL AND owner_id = ?", user.ID).First(pRootDir)
	return pRootDir, res.Error
}

func (user *User) Create() error {
	res := db.GetDB().Create(user)
	return res.Error
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	res := tx.Create(&Directory{
		Name:    "/",
		OwnerID: user.ID,
	})
	return res.Error
}
