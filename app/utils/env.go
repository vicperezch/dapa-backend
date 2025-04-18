package utils

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	envLoaded bool
)

// Carga el archivo .env una sola vez
func Load() {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("No se pudo cargar el archivo .env: %v", err)
		}
		envLoaded = true
	})
}

// Obtiene una variable de entorno o devuelve un valor por defecto 
func EnvGet(key, defaultValue string) string {
	//Si el archivo .env no esta cargado lo carga
	if !envLoaded {
		Load()
	}
	
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	
	if defaultValue != "" {
		return defaultValue
	}
	
	log.Printf("Variable de entorno %s no definida y sin valor por defecto", key)
	return ""
}

// Obtiene una variable de entorno o panic si no está definida
func EnvMustGet(key string) string {
	//Si el archivo .env no esta cargado lo carga
	if !envLoaded {
		Load()
	}
	
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf("Variable de entorno requerida %s no está definida", key)
	}
	
	return value
}