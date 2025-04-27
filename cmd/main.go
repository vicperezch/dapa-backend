package main

import (
	"log"

	"dapa/app/utils"
	"dapa/app/routes"
	"dapa/database"
	"dapa/app/model"
	_ "dapa/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			De aquí para allá API
//	@version		0.1
//	@description	API that provides the backend for the DAPA page.
//
//	@host			localhost:8080
//	@BasePath		/api
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

	database.ConnectToDatabase()
	if err := CreateFirstAdmin(); err != nil {
		log.Fatal("Error creando admin inicial: ", err)
	}

	// Configurar rutas
	routes.SetupRoutes(router)

	// Endpoint para documentación
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Iniciar servidor
	router.Run(":8080")
}

func CreateFirstAdmin() error {
	// Verificar si ya existe un admin
	var count int64
	db := database.DB
	db.Model(&model.Employee{}).Where("role = ?", "admin").Count(&count)
	
	if count > 0 {
		return nil // Ya existe admin
	}

	hashedPassword, err := utils.HashPassword("dapa12345")
	if err != nil {
		return err
	}

	admin := model.Employee{
		User: model.User{
			Name:     "Admin",
			Email:    "admin@dapa.com",
			Phone:    "0000000000",
		},
		Password: hashedPassword,
		Role:     "admin",
	}

	return db.Create(&admin).Error
}
