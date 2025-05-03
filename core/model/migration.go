package model

import (
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(&Directory{}, &File{}, &User{})
}
