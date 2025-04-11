package main

import (
	"dapa/app/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configuración CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Ping endpoint para probar conexión
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "¡Backend conectado con éxito!",
		})
	})
	

	// Endpoints de API reales
	api := router.Group("/api")
	{
		api.POST("/user", handlers.CreateUser)
		api.GET("/user", handlers.GetUsers)
	}

	// Iniciar servidor
	router.Run(":8080")
}
