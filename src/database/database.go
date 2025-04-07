package database

import (
	"log"
  "os"

  "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDataBase() *gorm.DB  {
  env := godotenv.Load(".env")
  if env != nil {
    panic("Error al cargar el archivo .env")
}
  POSTGRES_USER := os.Getenv("POSTGRES_USER")
  POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
  POSTGRES_DB := os.Getenv("POSTGRES_DB")

  dsn := "host=database user=" +POSTGRES_USER+" password="+POSTGRES_PASSWORD+" dbname="+POSTGRES_DB+" port=5432 sslmode=disable TimeZone=UTC"
  var err error

  db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
    log.Fatal("Failed to connect to DB", err)
  }

  return db
}