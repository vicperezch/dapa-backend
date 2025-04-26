package routes

import (
	"dapa/app/handlers"
	"dapa/app/middlewares"

	"github.com/gin-gonic/gin"
)

//Se encarga de manejar todas las rutas que admita el API
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

			// Rutas que pueden ser accedidas por cualquier rol
			protected.PUT("/users/:id", handlers. UpdateUser)
			
			// Rutas que solo el rol admin puede tener
			admin := protected.Group("")
			admin.Use(middlewares.RoleRequired("admin"))
			{
				admin.POST("/users", handlers.RegisterHandler)
				admin.GET("/users", handlers.GetUsers)
				admin.GET("/users/:id", handlers.GetUserById)
			}
			
			// Rutas que solo el rol driver puede tener
			driver := protected.Group("")
			driver.Use(middlewares.RoleRequired("driver"))
			{
				 
			}
		}
	}
}