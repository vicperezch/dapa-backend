package database

import (
	"log"

	"dapa/app/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	migrateDatabase()
}

// Aplica los scripts sql necesarios para migrar la base de datos
func migrateDatabase() {
	dbURL := 
		"postgres://" + utils.EnvMustGet("POSTGRES_USER") + 
		":" + utils.EnvMustGet("POSTGRES_PASSWORD") + 
		"@database:5432/" + utils.EnvMustGet("POSTGRES_DB") + 
		"?sslmode=disable"

	m, err := migrate.New(
		"file:///database/migrations",
		dbURL,		
	)

	if err != nil {
		log.Fatal("Failed to mirate DB", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to migrate DB", err)
	}
}
