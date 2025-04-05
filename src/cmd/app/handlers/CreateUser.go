package handlers

import (
	"log"
	"net/http"
	"dapa/src/cmd/app/models"
	"dapa/src/cmd/app/utils"
	"dapa/src/cmd/app/database"

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
		"INSERT INTO users (name, lastName, phone, email) VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), DEFAULT))",
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