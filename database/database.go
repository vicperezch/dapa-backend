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

// ConnectToDatabase initializes the connection to the PostgreSQL database using GORM,
// then applies any pending migrations.
func ConnectToDatabase() {
	dsn := "host=database user=" + utils.EnvMustGet("POSTGRES_USER") +
		" password=" + utils.EnvMustGet("POSTGRES_PASSWORD") +
		" dbname=" + utils.EnvMustGet("POSTGRES_DB") + " port=5432 sslmode=disable TimeZone=UTC"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	migrateDatabase()
}

// migrateDatabase runs SQL migration scripts from the migrations folder
// to keep the database schema up to date.
func migrateDatabase() {
	dbURL := "postgres://" + utils.EnvMustGet("POSTGRES_USER") +
		":" + utils.EnvMustGet("POSTGRES_PASSWORD") +
		"@database:5432/" + utils.EnvMustGet("POSTGRES_DB") +
		"?sslmode=disable"

	m, err := migrate.New(
		"file:///database/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to migrate DB:", err)
	}
}