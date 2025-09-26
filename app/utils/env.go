package utils

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once      sync.Once
	envLoaded bool
)

// Carga el archivo .env y sus variables
func Load() {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v", err)
		}
		envLoaded = true
	})
}

// Obtiene el valor de una variable del archivo .env
// Recibe la llave (nombre de la variable) y el valor por defecto a tomar
// Retorna el valor encontrado si la variable existe
// Si no, retorna el valor por defecto
func EnvGet(key, defaultValue string) string {
	if !envLoaded {
		Load()
	}

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue != "" {
		return defaultValue
	}

	log.Printf("Environment variable %s not set and no default value provided", key)
	return ""
}

// Obtiene el valor de una variable del archivo .env
// Si la variable no existe, detiene la ejecución del programa
func EnvMustGet(key string) string {
	if !envLoaded {
		Load()
	}

	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf("Required environment variable %s is not defined", key)
	}

	return value
}
