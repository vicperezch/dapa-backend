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

	api.GET("form/questions", handlers.GetQuestionsHandler)
	api.POST("form/submissions", handlers.CreateSubmissionHandler)

	// Protected routes - require valid JWT token
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())

	// Routes accessible to all authenticated users regardless of role
	{
		protected.PUT("/users/:id", handlers.UpdateUserHandler)
		protected.GET("/users/:id", handlers.GetUserHandler)

		protected.GET("/orders", handlers.GetOrdersHandler)
		protected.GET("/orders/:id", handlers.GetOrderHandler)
		protected.PATCH("/orders/:id/status", handlers.ChangeOrderStatusHandler)

		protected.GET("form/submissions/:id", handlers.GetSubmissionHandler)
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
		admin.POST("/orders", handlers.CreateOrderHandler)
		admin.PUT("/orders/:id", handlers.UpdateOrderHandler)
		admin.PATCH("/orders/:id/assign", handlers.AssignOrderHandler)

		// FORMULARIO: Tipos de pregunta
		admin.GET("form/question-types", handlers.GetQuestionTypesHandler)

		// FORMULARIO: Preguntas
		admin.POST("form/questions", handlers.CreateQuestionHandler)
		admin.PUT("form/questions/:id", handlers.UpdateQuestionHandler)
		admin.DELETE("form/questions/:id", handlers.DeleteQuestionHandler)
		admin.PATCH("form/questions/reorder", handlers.ReorderQuestionsHandler)
		admin.PATCH("form/questions/:id/active", handlers.ToggleQuestionActiveHandler)

		// FORMULARIO: Opciones de pregunta
		admin.POST("form/questions/:questionId/options", handlers.CreateQuestionOptionHandler)

		// FORMULARIO: Envíos (admin puede ver y actualizar estado)
		admin.GET("form/submissions", handlers.GetSubmissionsHandler)
		admin.PATCH("form/submissions/:id/status", handlers.UpdateSubmissionStatusHandler)

		// REPORTE: Financiero (detallado)
		admin.GET("/reports/financial", handlers.FinancialReport)
		admin.GET("/reports/financial/date", handlers.FinancialReportByDate)
		admin.GET("/reports/drivers", handlers.DriversReport)
	}
}
