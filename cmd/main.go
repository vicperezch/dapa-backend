package main

import (
	"log"

	"dapa/app/model"
	"dapa/app/routes"
	"dapa/app/utils"
	"dapa/database"
	_ "dapa/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           De aquí para allá API
// @version         0.1
// @description     API that provides the backend for the DAPA page.
//
// @host            localhost:8080
// @BasePath        /api
func main() {
	// Create a new Gin router instance with default middleware (logger and recovery)
	router := gin.Default()

	// Configure CORS middleware to allow cross-origin requests
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register custom validators for request binding
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validrole", utils.RoleValidator)
		v.RegisterValidation("password", utils.PasswordValidator)
		v.RegisterValidation("phone", utils.PhoneValidator)
		v.RegisterValidation("question_text", utils.QuestionTextValidator)
		v.RegisterValidation("question_desc", utils.QuestionDescriptionValidator)
		v.RegisterValidation("question_type", utils.QuestionTypeValidator)
		v.RegisterValidation("question_option", utils.QuestionOptionValidator)
		v.RegisterValidation("submission_status", utils.SubmissionStatusValidator)
	}

	// Connect to the database
	database.ConnectToDatabase()

	// Create initial admin user if none exists
	if err := CreateFirstAdmin(); err != nil {
		log.Fatal("Error creating initial admin: ", err)
	}

	// Setup all routes for the API
	routes.SetupRoutes(router)

	// Setup Swagger endpoint for API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server on port 8080
	router.Run(":8080")
}

// CreateFirstAdmin checks if an admin user exists and creates one if not
func CreateFirstAdmin() error {
	var count int64
	db := database.DB

	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("dapa12345")
	if err != nil {
		return err
	}

	// Prepare admin user data
	admin := model.User{
		Name:         "Admin",
		Email:        "admin@dapa.com",
		Phone:        "0000000000",
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	// Insert the new admin into the database
	return db.Create(&admin).Error
}
