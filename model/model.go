package model

import (
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initializes the database
func Setup() {
	var err error
	dsn := "sqlserver://sa:12345678@localhost:1433?database=godb"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	db.AutoMigrate(&User{})
}
