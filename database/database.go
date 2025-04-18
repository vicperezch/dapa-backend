package database

import (
	"log"

	"dapa/app/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDatabase() *gorm.DB {
	dsn := 
		"host=database user=" + utils.EnvMustGet("POSTGRES_USER") + 
		" password=" + utils.EnvMustGet("POSTGRES_PASSWORD") + 
		" dbname=" + utils.EnvMustGet("POSTGRES_DB") + " port=5432 sslmode=disable TimeZone=UTC"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	return db
}
