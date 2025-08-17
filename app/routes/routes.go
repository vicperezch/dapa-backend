package routes

import (
	"dapa/app/handlers"
	"dapa/app/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and applies middleware for authentication and authorization.
func SetupRoutes(router *gin.Engine) {

	// Public routes - no authentication required
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Backend connected successfully!",
		})
	})

	// Base API group
	api := router.Group("/api")

	// Public authentication route
	api.POST("/login", handlers.LoginHandler)

	// Protected routes - require valid JWT token
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	// Routes accessible to all authenticated users regardless of role
	{
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.GET("/users/:id", handlers.GetUserById)
	}

	// Admin-only routes group
	admin := protected.Group("")
	admin.Use(middlewares.RoleRequired("admin"))
	{
		// User management routes for admin
		admin.POST("/users", handlers.RegisterHandler)
		admin.GET("/users", handlers.GetUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// Vehicle management routes for admin
		admin.GET("/vehicles", handlers.GetVehicles)
		admin.POST("/vehicles", handlers.CreateVehicle)
		admin.GET("/vehicles/:id", handlers.GetVehicleById)
		admin.PUT("/vehicles/:id", handlers.UpdateVehicle)
		admin.DELETE("/vehicles/:id", handlers.DeleteVehicle)

		admin.GET("/quotes", handlers.GetQuotes)
		admin.POST("/quotes/:id", handlers.AssignQuoteInfo)

		// FORMULARIO: Tipos de pregunta
		admin.GET("/question-types", handlers.GetQuestionTypes)
		admin.POST("/question-types", handlers.CreateQuestionType)

		// FORMULARIO: Preguntas
		admin.GET("/questions", handlers.GetQuestions)
		admin.GET("/questions-active", handlers.GetActiveQuestions)
		admin.POST("/questions", handlers.CreateQuestion)
		admin.PUT("/questions/:id", handlers.UpdateQuestion)
		admin.DELETE("/questions/:id", handlers.DeleteQuestion)
		admin.PATCH("/questions/reorder", handlers.ReorderQuestions)
		admin.PATCH("/questions/:id/active", handlers.ToggleQuestionActive)

		// FORMULARIO: Opciones de pregunta
		admin.POST("/questions/:questionId/options", handlers.CreateQuestionOption)

		// FORMULARIO: Envíos (admin puede ver y actualizar estado)
		admin.GET("/submissions", handlers.GetSubmissions)
		admin.GET("/submissions-stats", handlers.GetSubmissionStats)
		admin.PUT("/submissions/:id/status", handlers.UpdateSubmissionStatus)
	}

	// Driver-only routes group (currently empty, placeholder for future driver endpoints)
	driver := protected.Group("")
	driver.Use(middlewares.RoleRequired("driver"))
	{
		// Add driver-specific routes here
	}
}

