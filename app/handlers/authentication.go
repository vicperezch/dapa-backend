package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterHandler(c *gin.Context) {
	db := database.ConnectToDatabase()

	var req model.RegisterRequest
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request: ", err)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make([]string, len(ve))

			for i, fe := range ve {
				errorMessages[i] = utils.GetTagMessage(fe.Tag())
			}

			utils.RespondWithError(c, errorMessages[0], http.StatusBadRequest)
			return
		}

		utils.RespondWithError(c, err.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password: ", err)
		utils.RespondWithError(c, "Error registering user", http.StatusInternalServerError)
		return
	}

	employee := model.Employee{
		User: model.User{
			Name:     req.Name,
			LastName: req.LastName,
			Phone:    req.Phone,
			Email:    req.Email,
		},
		Password: passwordHash,
		Role:     req.Role,
	}
	if err = db.Create(&employee).Error; err != nil {
		utils.RespondWithError(c, "Error registering user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "User created successfully",
	})
}

func LoginHandler(c *gin.Context) {
	db := database.ConnectToDatabase()

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request: ", err)
		utils.RespondWithError(c, "Invalid request", http.StatusBadRequest)
		return
	}

	var employee model.Employee

	err := db.Preload("User").
		Joins("JOIN users ON users.id = employees.user_id").
		Where("users.email = ?", req.Email).
		First(&employee).Error

	if err != nil {
		log.Println("Error finding user: ", err)
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(req.Password, employee.Password) {
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Login successful",
	})
}