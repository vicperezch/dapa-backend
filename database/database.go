package database

import (
	"log"

	"dapa/app/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := 
		"host=database user=" + utils.EnvMustGet("POSTGRES_USER") + 
		" password=" + utils.EnvMustGet("POSTGRES_PASSWORD") + 
		" dbname=" + utils.EnvMustGet("POSTGRES_DB") + " port=5432 sslmode=disable TimeZone=UTC"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}
}
