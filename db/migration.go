package db

import (
	"conviction/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(&model.User{}, &model.File{})
}
