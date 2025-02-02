package model

import (
	"conviction/db"

	"gorm.io/gorm"
)

type Directory struct {
	gorm.Model
	Name     string ``
	OwnerID  uint
	ParentID uint
}

func (d *Directory) GetChild(dirName string) (*Directory, error) {
	var childDir Directory
	res := db.GetDB().
		Where("parent_id = ? AND owner_id = ? AND name = ?", d.ID, d.OwnerID, dirName).
		First(&childDir)
	return &childDir, res.Error
}

// insert directory into DB table
func (d *Directory) Create() error {
	if res := db.GetDB().FirstOrCreate(d, *d); nil == res.Error {
		return nil
	}

	d.Model = gorm.Model{}

	res := db.GetDB().First(d)
	return res.Error
}

func (d *Directory) GetChildDirectory() ([]Directory, error) {
	var childDir []Directory
	res := db.GetDB().
		Where("parent_id = ?", d.ID).
		Find(childDir)
	return childDir, res.Error
}

func (d *Directory) GetChildFile() ([]File, error) {
	var childFile []File
	res := db.GetDB().
		Where("directory_id = ?", d.ID).
		Find(&childFile)
	return childFile, res.Error
}
