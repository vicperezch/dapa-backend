package handlers

import (
	"log"
	"net/http"
	"dapa/app/models"
	"dapa/app/utils"
	"dapa/database"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	db := database.ConnectToDataBase()
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		utils.RespondWithError(c, "Error getting all users", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, users)
}