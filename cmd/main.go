package main

import (
	"dapa/app/handlers"
	"dapa/app/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// Implementación de los custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validrole", utils.RoleValidator)
		v.RegisterValidation("password", utils.PasswordValidator)
		v.RegisterValidation("phone", utils.PhoneValidator)
	}

	// Ping endpoint para probar conexión
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "¡Backend conectado con éxito!",
		})
	})

	// Endpoints de API reales
	api := router.Group("/api")
	{
		api.POST("/users", handlers.RegisterHandler)
		api.GET("/users", handlers.GetUsers)
	}

	// Iniciar servidor
	router.Run(":8080")
}
