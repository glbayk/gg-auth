package models

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	// connect to database
	dsn := os.Getenv("POSTGRES_URL_GORM")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
		os.Exit(1)
	}

	// migrate database
	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error migrating database")
		os.Exit(1)
	}
}
