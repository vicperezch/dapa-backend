package routes

import (
	"dapa/app/handlers"
	"dapa/app/middlewares"

	"github.com/gin-gonic/gin"
)

// Se encarga de manejar todas las rutas que admita el API
func SetupRoutes(router *gin.Engine) {

	// Rutas públicas
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "¡Backend conectado con éxito!",
		})
	})

	api := router.Group("/api")
	{
		api.POST("/login", handlers.LoginHandler)

		// Rutas protegidas, solo pueden ser accedidas si se cuenta con un token
		protected := api.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Rutas que pueden ser accedidas por cualquier rol autenticado
			protected.PUT("/users/:id", handlers.UpdateUser)
			protected.GET("/users/:id", handlers.GetUserById)

			// FORMULARIO: Listar preguntas y tipos (para clientes y drivers)
			protected.GET("/question-types", handlers.GetQuestionTypes)
			protected.GET("/questions", handlers.GetQuestions)
			protected.GET("/questions/:id", handlers.GetQuestionByID)

			// FORMULARIO: Envíos (clientes pueden crear y ver sus submissions)
			protected.POST("/submissions", handlers.CreateSubmission)
			protected.GET("/submissions", handlers.GetSubmissions) // Puedes filtrar por usuario en el handler

			// Rutas que solo el rol admin puede tener
			admin := protected.Group("")
			admin.Use(middlewares.RoleRequired("admin"))
			{
				admin.POST("/users", handlers.RegisterHandler)
				admin.GET("/users", handlers.GetUsers)
				admin.DELETE("/users/:id", handlers.DeleteUser)

				admin.GET("/vehicles", handlers.GetVehicles)
				admin.POST("/vehicles", handlers.CreateVehicle)
				admin.GET("/vehicles/:id", handlers.GetVehicleById)
				admin.PUT("/vehicles/:id", handlers.UpdateVehicle)
				admin.DELETE("/vehicles/:id", handlers.DeleteVehicle)

				// FORMULARIO: Tipos de pregunta
				admin.POST("/question-types", handlers.CreateQuestionType)

				// FORMULARIO: Preguntas
				admin.POST("/questions", handlers.CreateQuestion)
				admin.PUT("/questions/:id", handlers.UpdateQuestion)
				admin.DELETE("/questions/:id", handlers.DeleteQuestion)

				// FORMULARIO: Opciones de pregunta
				admin.POST("/questions/:questionId/options", handlers.CreateQuestionOption)

				// FORMULARIO: Envíos (admin puede ver y actualizar estado)
				admin.PUT("/submissions/:id/status", handlers.UpdateSubmissionStatus)
			}

			// Rutas que solo el rol driver puede tener
			driver := protected.Group("")
			driver.Use(middlewares.RoleRequired("driver"))
			{
				// Ejemplo: El driver puede ver submissions asignados a su viaje
				driver.GET("/driver/submissions", handlers.GetDriverSubmissions)
			}
		}
	}
}
