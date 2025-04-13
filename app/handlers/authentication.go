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
