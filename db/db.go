package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var g_db *gorm.DB

func InitDB() {
	dsn := "host=127.0.0.1 user=postgres password=0403 dbname=convictiodb port=5432 sslmode=disable"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	g_db = db
}

func GetDB() *gorm.DB {
	return g_db
}
