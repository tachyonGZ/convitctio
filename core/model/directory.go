package model

import (
	"conviction/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Directory struct {
	UUID string `gorm:"column:uuid;primarykey"`

	OwnerUUID  string  `gorm:"column:owner_uuid;index:owner_uuid"`
	ParentUUID *string `gorm:"column:parent_uuid;index:parent_id;unique_index:idx_only_one_name"`

	Name string `gorm:"unique_index:idx_only_one_name"`
}

func (pDir *Directory) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = errors.New("uuid grenate failed")
	}
	pDir.UUID = "dir_" + uuid.String()
	return
}

func (d *Directory) BeforeDelete(tx *gorm.DB) (err error) {
	//  if u.Confirmed {
	//    tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
	//  }

	childDir := Directory{}
	resDir := db.GetDB().Where("parent_uuid = ?", d.UUID).Find(&childDir)
	if resDir.Error != nil {
		err = resDir.Error
	}

	childFile := File{}
	resFile := db.GetDB().Where("directory_uuid = ?", d.UUID).Find(&childFile)
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
		Where("parent_uuid = ? AND owner_uuid = ? AND name = ?", d.UUID, d.OwnerUUID, dirName).
		Find(pChildDir)
	return pChildDir, res.RowsAffected != 0, res.Error
}

func (pDir *Directory) Create() error {
	res := db.GetDB().Create(pDir)
	return res.Error
}

func (d *Directory) GetChildDirectory() (childDir []Directory, err error) {
	res := db.GetDB().Where("parent_uuid = ?", d.UUID).Find(&childDir)
	err = res.Error
	return
}

func (d *Directory) GetChildFile() (childFile []File, err error) {
	res := db.GetDB().Where("directory_uuid = ?", d.UUID).Find(&childFile)
	err = res.Error
	return
}

func FindUserDirectory(user_uuid string, dir_uuid string) (*Directory, error) {
	dir := Directory{}
	res := db.GetDB().
		Where("owner_uuid = ? AND uuid = ?", user_uuid, dir_uuid).
		Find(&dir)
	return &dir, res.Error
}

func DeleteUserDirectory(user_uuid string, dir_uuid string) error {
	res := db.GetDB().
		Unscoped().
		Where("owner_uuid = ? AND uuid = ?", user_uuid, dir_uuid).
		Delete(&Directory{})
	return res.Error
}

func GetUserRootID(user_uuid string) (string, error) {
	root_dir := Directory{}
	res := db.GetDB().
		Where("owner_uuid = ? AND parent_uuid is NULL", user_uuid).
		Find(&root_dir)
	return root_dir.UUID, res.Error
}
