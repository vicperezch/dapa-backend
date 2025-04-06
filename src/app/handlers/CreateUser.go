package handlers

import (
	"log"
	"net/http"
	"dapa/app/models"
	"dapa/app/utils"
	"dapa/database"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
}

func CreateUser(c *gin.Context) {
	db := database.ConnectToDataBase()
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	result := db.Exec(
		"INSERT INTO users (name, last_name, phone, email) VALUES ($1, $2, $3, NULLIF($4, ''))",
		req.Name, req.LastName, req.Phone, req.Email,
	)

	if result.Error != nil {
		log.Println("Error creating new user:", result.Error)
		utils.RespondWithError(c, "Error creating new user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, models.ApiResponse{
		Success: true,
		Message: "Successfully created user",
	})
}