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

	CreatorID string `gorm:"column:creator_uuid"` // user id
	SourceID  string `gorm:"column:source_uuid"`  // file id
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
