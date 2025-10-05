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
	router := gin.Default()

	// Configuración del middleware de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Validaciones personalizadas
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", utils.PasswordValidator)
		v.RegisterValidation("phone", utils.PhoneValidator)
		v.RegisterValidation("plate", utils.LicensePlateValidator)
		v.RegisterValidation("question_text", utils.QuestionTextValidator)
		v.RegisterValidation("question_desc", utils.QuestionDescriptionValidator)
		v.RegisterValidation("question_type", utils.QuestionTypeValidator)
		v.RegisterValidation("question_option", utils.QuestionOptionValidator)
		v.RegisterValidation("submission_status", utils.SubmissionStatusValidator)
	}

	database.ConnectToDatabase()

	// Crear usuario administrador de prueba
	if err := CreateFirstAdmin(); err != nil {
		log.Fatal("Error creating initial admin: ", err)
	}
	SeedQuestionTypes()

	routes.SetupRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}

func CreateFirstAdmin() error {
	var count int64
	db := database.DB

	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashString("dapa12345")
	if err != nil {
		return err
	}

	admin := model.User{
		Name:         "Admin",
		Email:        "admin@dapa.com",
		Phone:        "0000000000",
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	return db.Create(&admin).Error
}

func SeedQuestionTypes() {
	questionTypes := []string{"text", "multiple", "unique", "dropdown", "area"}

	for _, typeName := range questionTypes {
		var existingType model.QuestionType
		result := database.DB.Where("type = ?", typeName).First(&existingType)

		if result.Error != nil {
			newType := model.QuestionType{Type: typeName}
			if err := database.DB.Create(&newType).Error; err != nil {
				log.Printf("Error creating question type %s: %v", typeName, err)
			} else {
				log.Printf("Question type '%s' created with ID: %d", typeName, newType.ID)
			}
		} else {
			log.Printf("Question type '%s' already exists with ID: %d", typeName, existingType.ID)
		}
	}
}
