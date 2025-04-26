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
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	
	if claims.Role != "admin" {
		utils.RespondWithError(c,"Insufficient permissions",http.StatusForbidden )
		return
	}

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

func GetUserById(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c,"Insufficient permissions",http.StatusForbidden )
		return
	}

	db := database.ConnectToDatabase()
	var user model.User

	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		log.Println("Error fetching user:", err)
		utils.RespondWithError(c, "Error getting user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, user)
}

func UpdateUser(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	db := database.ConnectToDatabase()
	var req model.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	id := c.Param("id")

	if claims.Role == "admin" {
		result := db.Exec(
			"UPDATE users SET name = $1, last_name = $2, phone = $3, email = NULLIF($4, '') WHERE id = $5",
			req.Name, req.LastName, req.Phone, req.Email, id,
		)
	
		if result.Error != nil {
			log.Println("Error updating user:", result.Error)
			utils.RespondWithError(c, "Error updating user", http.StatusInternalServerError)
			return
		}
	
		result = db.Exec(
			"UPDATE employees SET role = $1 WHERE user_id = $2",
			req.Role, id,
		)
	
		if result.Error != nil {
			log.Println("Error updating employee role:", result.Error)
			utils.RespondWithError(c, "Error updating user role", http.StatusInternalServerError)
			return
		}
	
	} else {
		result := db.Exec(
			"UPDATE users SET name = $1, last_name = $2, phone = $3, email = NULLIF($4, '') WHERE id = $5",
			req.Name, req.LastName, req.Phone, req.Email, id,
		)
	
		if result.Error != nil {
			log.Println("Error updating user:", result.Error)
			utils.RespondWithError(c, "Error updating user", http.StatusInternalServerError)
			return
		}
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully updated user",
	})
}