package model

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	UserID uint   `gorm:"index:user_id;unique_index:idx_only_one"`
	Name   string `gorm:"unique_index:idx_only_one"`
	Path   string `gorm:"type:text"`
}

func GetFileByID(ID uint) File {
	var u File
	return u
}

func (f *File) Insert() {

}
