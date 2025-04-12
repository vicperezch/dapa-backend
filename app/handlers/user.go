package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	db := database.ConnectToDatabase()
	var users []model.User

	if err := db.Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		utils.RespondWithError(c, "Error getting all users", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, users)
}

func CreateUser(c *gin.Context) {
	db := database.ConnectToDatabase()
	var req model.CreateUserRequest

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

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully created user",
	})
}
