package models

import (
	"xiaoshuo-backend/config"
	"gorm.io/gorm"
)

// DB is the database instance from config
var DB *gorm.DB

// InitializeDB sets the database instance from config
func InitializeDB() {
	DB = config.DB
}