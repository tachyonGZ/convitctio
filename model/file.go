package model

import (
	"conviction/db"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	UserID      uint   `gorm:"index:user_id;unique_index:idx_only_one"`
	Name        string `gorm:"unique_index:idx_only_one"`
	Path        string `gorm:"type:text"`
	Size        uint64
	DirectoryID uint `gorm:"index:directory_id;unique_index:idx_only_one"`
}

func IsSameNameFileExist(name string, dirID uint, userID uint) bool {
	file := &File{}
	res := db.GetDB().Where("name = ? AND directory_id = ? AND user_id = ?", name, dirID, userID).Find(file)
	return res.RowsAffected != 0
}

func (file *File) Create() error {
	res := db.GetDB().Create(file)
	return res.Error
}

func (file *File) Delete() error {
	res := db.GetDB().Delete(file)
	return res.Error
}

func DeleteUserFile(userID uint, fileID string) error {
	res := db.GetDB().Unscoped().Where("user_id = ? AND id = ?", userID, fileID).Delete(&File{})
	return res.Error
}

func GetFileByID(fileID uint, userID uint) (*File, error) {
	f := File{}
	res := db.GetDB().Where("id = ? AND user_id = ?", fileID, userID).First(&f)
	return &f, res.Error
}

func FindUserFile(userID uint, fileID string) (*File, error) {
	file := File{}
	res := db.GetDB().Where("user_id = ? AND id = ?", userID, fileID).Find(&file)
	return &file, res.Error
}

func (file *File) PlaceholderToFile() {

}
