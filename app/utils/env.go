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

// Load reads the .env file once and loads environment variables into the process.
// It uses sync.Once to ensure the file is loaded only a single time during the app lifecycle.
func Load() {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v", err)
		}
		envLoaded = true
	})
}

// EnvGet retrieves the value of the environment variable named by the key.
// If the variable is not set, it returns the provided defaultValue.
// If no defaultValue is provided and the variable is not set, it logs a warning and returns an empty string.
// It ensures the .env file is loaded before accessing environment variables.
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

// EnvMustGet retrieves the value of the environment variable named by the key.
// If the variable is not set, it logs a panic and stops execution.
// It ensures the .env file is loaded before accessing environment variables.
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