package routes

import (
	"dapa/app/handlers"
	"dapa/app/middlewares"

	"github.com/gin-gonic/gin"
)

// Configura todal las rutas de la API y aplica los middleware necesarios
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas públicas
	// Autenticación
	api.POST("/login", handlers.LoginHandler)
	api.POST("/auth/forgot", handlers.ForgotPasswordHandler)
	api.POST("/auth/reset", handlers.ResetPasswordHandler)

	// Formulario para clientes
	api.GET("form/questions", handlers.GetQuestionsHandler)
	api.POST("form/submissions", handlers.CreateSubmissionHandler)

	// Tracking de órdenes
	api.GET("/orders/track", handlers.OrderTrackingHandler)

	// Rutas que requieren que el usuario se encuentra autenticado
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())
	{
		// ENTIDADES: Usuarios
		protected.PUT("/users/:id", handlers.UpdateUserHandler)
		protected.GET("/users/:id", handlers.GetUserHandler)

		// ENTIDADES: órdenes
		protected.GET("/orders/:id/token", handlers.GetOrderTokenHandler)  // MÁS ESPECÍFICO PRIMERO
		protected.GET("/orders/:id", handlers.GetOrderHandler)             // MENOS ESPECÍFICO DESPUÉS
		protected.GET("/orders", handlers.GetOrdersHandler)
		protected.PATCH("/orders/:id/status", handlers.ChangeOrderStatusHandler)

		// FORMULARIO: respuestas
		protected.GET("form/submissions/:id", handlers.GetSubmissionHandler)
	}

	// Rutas que requieren que el usuario posea el rol de admin
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
		admin.PATCH("form/questions/:id/required", handlers.ToggleQuestionRequiredHandler)

		// FORMULARIO: Opciones de pregunta
		admin.POST("form/questions/:questionId/options", handlers.CreateQuestionOptionHandler)

		// FORMULARIO: Envíos (admin puede ver y actualizar estado)
		admin.GET("form/submissions", handlers.GetSubmissionsHandler)
		admin.PATCH("form/submissions/:id/status", handlers.UpdateSubmissionStatusHandler)

		// REPORTE: Financiero (detallado)
		admin.GET("/reports/financial", handlers.FinancialReport)
		admin.GET("/reports/financial/date", handlers.FinancialReportByDate)
		admin.GET("/reports/drivers", handlers.DriversReport)
		admin.GET("/reports/income", handlers.TotalIncomeReport)
		admin.GET("/reports/financial-control-income", handlers.FinancialControlIncome)
		admin.GET("/reports/financial-control-spending", handlers.FinancialControlSpending)

		// REPORTE: Gráficas
		admin.GET("/reports/completed-quotations", handlers.CompletedQuotationsChart)
		admin.GET("/reports/quotations-status", handlers.QuotationsStatusChart)
		admin.GET("/reports/drivers-performance", handlers.DriversPerformanceChart)
		admin.GET("/reports/drivers-trip-participation", handlers.DriversTripParticipationChart)
	}
}
