package db

import (
	"conviction/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

func Init(host string, user string, password string, name string, port uint16) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, name, port)

	gorm_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if config.SystemConfig.Debug {
		gorm_db.Logger.LogMode(logger.Info)
	} else {
		gorm_db.Logger.LogMode(logger.Silent)
	}

	GormDB = gorm_db
}
