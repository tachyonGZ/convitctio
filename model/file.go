package model

import (
	"conviction/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	UUID string `gorm:"column:uuid"`

	UserID      uint   `gorm:"index:user_id;unique_index:idx_only_one"`
	OwnerUUID   string `gorm:"column:owner_uuid"`
	Name        string `gorm:"unique_index:idx_only_one"`
	Path        string `gorm:"type:text"`
	Size        uint64
	DirectoryID uint `gorm:"index:directory_id;unique_index:idx_only_one"`
}

func (pFile *File) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = errors.New("uuid grenate failed")
	}
	pFile.UUID = "file_" + uuid.String()
	return
}

func IsSameNameFileExist(name string, dirID uint, userID uint) bool {
	file := &File{}
	res := db.GetDB().Where("name = ? AND directory_id = ? AND user_id = ?", name, dirID, userID).Find(file)
	return res.RowsAffected != 0
}

func (pFile *File) Create() error {
	res := db.GetDB().Create(pFile)
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

func FindUserFile2(owner_uuid string, file_uuid string) (*File, error) {
	file := File{}
	res := db.GetDB().Where("owner_uuid = ? AND uuid = ?", owner_uuid, file_uuid).Find(&file)
	return &file, res.Error
}

func IsUserOwnFile(userID string, fileID string) (bool, error) {
	file := File{}
	res := db.GetDB().Where("owner_uuid = ? AND uuid = ?", userID, fileID).Find(&file)
	return res.RowsAffected != 0, res.Error
}

func (file *File) PlaceholderToFile() {

}
