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
	api.POST("/auth/forgot", handlers.ForgotPasswordHandler)
	api.POST("/auth/reset", handlers.ResetPasswordHandler)

	// Protected routes - require valid JWT token
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	// Routes accessible to all authenticated users regardless of role
	{
		protected.PUT("/users/:id", handlers.UpdateUserHandler)
		protected.GET("/users/:id", handlers.GetUserHandler)

		protected.GET("/orders", handlers.GetOrdersHandler)
	}

	admin := protected.Group("")
	admin.Use(middlewares.RoleRequired("admin"))
	{
		// ENTIDADES: Usuarios
		admin.POST("/users", handlers.RegisterHandler)
		admin.GET("/users", handlers.GetUsersHandler)
		admin.DELETE("/users/:id", handlers.DeleteUserHandler)

		// ENTIDADES: Vehículos
		admin.GET("/vehicles", handlers.GetVehiclesHandler)
		admin.POST("/vehicles", handlers.CreateVehicleHandler)
		admin.GET("/vehicles/:id", handlers.GetVehicleHandler)
		admin.PUT("/vehicles/:id", handlers.UpdateVehicleHandler)
		admin.DELETE("/vehicles/:id", handlers.DeleteVehicleHandler)

		// ENTIDADES: Órdenes
		admin.GET("/orders/:id", handlers.GetOrderHandler)
		admin.PUT("/orders/:id", handlers.UpdateOrderHandler)
		admin.PATCH("/orders/:id", handlers.AssignOrderHandler)

		// FORMULARIO: Tipos de pregunta
		admin.GET("/question-types", handlers.GetQuestionTypesHandler)

		// FORMULARIO: Preguntas
		admin.GET("/questions", handlers.GetQuestionsHandler)
		admin.POST("/questions", handlers.CreateQuestionHandler)
		admin.PUT("/questions/:id", handlers.UpdateQuestionHandler)
		admin.DELETE("/questions/:id", handlers.DeleteQuestionHandler)
		admin.PATCH("/questions/reorder", handlers.ReorderQuestionsHandler)
		admin.PATCH("/questions/:id/active", handlers.ToggleQuestionActiveHandler)

		// FORMULARIO: Opciones de pregunta
		admin.POST("/questions/:questionId/options", handlers.CreateQuestionOptionHandler)

		// FORMULARIO: Envíos (admin puede ver y actualizar estado)
		admin.GET("/submissions", handlers.GetSubmissionsHandler)
		admin.POST("/submissions", handlers.CreateSubmissionHandler)
		admin.GET("/submissions/:id", handlers.GetSubmissionHandler)
		admin.PUT("/submissions/:id/status", handlers.UpdateSubmissionStatusHandler)
	}
}
