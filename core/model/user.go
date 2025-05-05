package model

import (
	"conviction/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID string `gorm:"column:uuid;primarykey"`

	Username string `gorm:"size:50"`
	Password string `json:"-"`
}

func (pUser *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = errors.New("uuid grenate failed")
	}
	pUser.UUID = "user_" + uuid.String()
	return
}

func FindUser(user_uuid string) (*User, error) {
	user := User{}
	res := db.GormDB.Where("uuid", user_uuid).Find(&user)
	return &user, res.Error
}

func FindUserByUsername(username string) (User, error) {
	var user User
	result := db.GormDB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func (user *User) CheckPassword(password string) bool {
	return password == user.Password
}

func (pUser *User) Create() error {
	res := db.GormDB.Create(pUser)
	return res.Error
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	res := tx.Create(&Directory{
		Name:      "/",
		OwnerUUID: user.UUID,
	})
	return res.Error
}

func (user *User) Root() (*Directory, error) {
	pRootDir := &Directory{}
	res := db.GormDB.Where("parent_uuid is NULL AND owner_uuid = ?", user.UUID).First(pRootDir)
	return pRootDir, res.Error
}
