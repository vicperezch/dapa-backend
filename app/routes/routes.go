package routes

import (
	"dapa/app/handlers"
	"dapa/app/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and applies middleware for authentication and authorization.
func SetupRoutes(router *gin.Engine) {
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

		protected.GET("/orders", handlers.GetOrders)
	}

	admin := protected.Group("")
	admin.Use(middlewares.RoleRequired("admin"))
	{
		// ENTIDADES: Usuarios
		admin.POST("/users", handlers.RegisterHandler)
		admin.GET("/users", handlers.GetUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// ENTIDADES: Vehículos
		admin.GET("/vehicles", handlers.GetVehicles)
		admin.POST("/vehicles", handlers.CreateVehicle)
		admin.GET("/vehicles/:id", handlers.GetVehicleById)
		admin.PUT("/vehicles/:id", handlers.UpdateVehicle)
		admin.DELETE("/vehicles/:id", handlers.DeleteVehicle)

		// ENTIDADES: Órdenes
		admin.GET("/orders/:id", handlers.GetOrderById)
		admin.PUT("/orders/:id", handlers.UpdateOrder)
		admin.PATCH("/orders/:id", handlers.AssignOrder)

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
		admin.POST("/submissions", handlers.CreateSubmission)
		admin.GET("/submissions/:id", handlers.GetSubmissionByID)
		admin.GET("/submissions-stats", handlers.GetSubmissionStats)
		admin.PUT("/submissions/:id/status", handlers.UpdateSubmissionStatus)
	}
}
