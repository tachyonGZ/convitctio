package model

import (
	"conviction/db"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SharedFile struct {
	gorm.Model
	UUID string `gorm:"column:uuid"`

	CreatorUUID string `gorm:"column:creator_uuid"` // user id
	SourceUUID  string `gorm:"column:source_uuid"`  // file id
}

func (pShare *SharedFile) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		err = errors.New("uuid grenate failed")
	}
	pShare.UUID = "sharedfile_" + uuid.String()
	return
}

func (pShare *SharedFile) Create() error {
	res := db.GetDB().Create(pShare)
	return res.Error
}

func FindSharedFile(shared_file_uuid string) (*SharedFile, error) {
	p_shared_file := &SharedFile{}
	res := db.GetDB().Where("uuid = ?", shared_file_uuid).Find(p_shared_file)
	return p_shared_file, res.Error
}
