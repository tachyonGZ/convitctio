package model

import (
	"conviction/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Directory struct {
	gorm.Model
	DirectoryUUID string

	Name     string `gorm:"unique_index:idx_only_one_name"`
	OwnerID  uint   `gorm:"index:owner_id"`
	ParentID *uint  `gorm:"index:parent_id;unique_index:idx_only_one_name"`
}

func (pDir *Directory) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = errors.New("uuid grenate failed")
	}
	pDir.DirectoryUUID = "dir_" + uuid.String()
	return
}

func (d *Directory) BeforeDelete(tx *gorm.DB) (err error) {
	//  if u.Confirmed {
	//    tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	//  }

	childDir := Directory{}
	resDir := db.GetDB().Where("parent_id = ?", d.ID).Find(&childDir)
	if resDir.Error != nil {
		err = resDir.Error
	}

	childFile := File{}
	resFile := db.GetDB().Where("directory_id = ?", d.ID).Find(&childFile)
	if resFile.Error != nil {
		err = resFile.Error
	}

	if resDir.RowsAffected != 0 || resFile.RowsAffected != 0 {
		err = errors.New("directory not empty")
	}

	return
}

func (d *Directory) GetChild(dirName string) (*Directory, bool, error) {
	pChildDir := &Directory{}
	res := db.GetDB().
		Where("parent_id = ? AND owner_id = ? AND name = ?", d.ID, d.OwnerID, dirName).
		Find(pChildDir)
	return pChildDir, res.RowsAffected != 0, res.Error
}

func (pDir *Directory) Create() error {
	res := db.GetDB().Create(pDir)
	return res.Error
}

func (d *Directory) GetChildDirectory() (childDir []Directory, err error) {
	res := db.GetDB().Where("parent_id = ?", d.ID).Find(&childDir)
	err = res.Error
	return
}

func (d *Directory) GetChildFile() (childFile []File, err error) {
	res := db.GetDB().Where("directory_id = ?", d.ID).Find(&childFile)
	err = res.Error
	return
}

func GetUserDirectory(userID uint, dirID uint) (*Directory, error) {
	dir := Directory{}
	res := db.GetDB().Where("owner_id = ? AND id = ?", userID, dirID).Find(&dir)
	return &dir, res.Error
}

func DeleteUserDirectory(userID uint, dirID string) error {
	res := db.GetDB().Unscoped().Where("owner_id = ? AND id = ?", userID, dirID).Delete(&Directory{})
	return res.Error
}

func GetUserRootID(userID uint) (uint, error) {
	rootDir := Directory{}
	res := db.GetDB().Where("parent_id is NULL AND owner_id = ?", userID).First(&rootDir)
	return rootDir.ID, res.Error
}
