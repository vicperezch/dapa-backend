package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"time"
)

// @Summary		Get all users
// @Description	Returns a list of all users in the system.
// @Tags		users
// @Produce		json
// @Success		200	{array} model.User "List of users"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching users"
// @Router		/users/ [get]
func GetUsers(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var users []model.User

	if err := database.DB.Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		utils.RespondWithError(c, "Error getting all users", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, users)
}

// @Summary		Get user by ID
// @Description	Returns the user information based on the given ID.
// @Tags		users
// @Produce		json
// @Param		id path int true "User ID"
// @Success		200	{object} model.User "User found"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error fetching user"
// @Router		/users/{id} [get]
func GetUserById(c *gin.Context) {

	var user model.User

	id := c.Param("id")
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Println("Error fetching user:", err)
		utils.RespondWithError(c, "Error getting user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, user)
}

// @Summary		Create a new user
// @Description	Creates a new user entry in the database.
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body model.CreateUserRequest true "User information to create"
// @Success		200	{object} model.ApiResponse "Successfully created user"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error creating new user"
// @Router		/users/ [post]
func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:     req.Name,
		LastName: req.LastName,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("Error creating new user:", err)
		utils.RespondWithError(c, "Error creating new user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully created user",
	})
}

// @Summary		Update user by ID
// @Description	Updates the user information and role based on the given ID.
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		id path int true "User ID"
// @Param		user body model.UpdateUserRequest true "Updated user information"
// @Success		200	{object} model.ApiResponse "Successfully updated user"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		500	{object} model.ApiResponse "Error updating user"
// @Router		/users/{id} [put]
func UpdateUser(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request:", err)
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	id := c.Param("id")

	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		log.Println("Error finding user:", err)
		utils.RespondWithError(c, "User not found", http.StatusInternalServerError)
		return
	}

	user.Name = req.Name
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Email = req.Email
	user.LastModifiedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		log.Println("Error updating user:", err)
		utils.RespondWithError(c, "Error updating user", http.StatusInternalServerError)
		return
	}

	if claims.Role == "admin" {
		if err := database.DB.Model(&model.Employee{}).Where("user_id = ?", id).Update("role", req.Role).Error; err != nil {
			log.Println("Error updating employee role:", err)
			utils.RespondWithError(c, "Error updating user role", http.StatusInternalServerError)
			return
		}
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully updated user",
	})
}

// @Summary		Mark user as inactive
// @Description	Marks the user as inactive instead of permanently deleting.
// @Tags		users
// @Produce		json
// @Param		id path int true "User ID"
// @Success		200	{object} model.ApiResponse "Successfully marked user as inactive"
// @Failure		403	{object} model.ApiResponse "Insufficient permissions"
// @Failure		500	{object} model.ApiResponse "Error updating user status"
// @Router		/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	id := c.Param("id")

	if err := database.DB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error; err != nil {
		log.Println("Error updating user to inactive:", err)
		utils.RespondWithError(c, "Error deleting user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Successfully marked user as inactive",
	})
}
