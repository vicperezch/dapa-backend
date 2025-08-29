package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
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

	if err := database.DB.Where("is_active = ?", true).Find(&users).Error; err != nil {
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
	if err := database.DB.
		Where("id = ? AND is_active = ?", id, true).
		First(&user).Error; err != nil {

		utils.RespondWithError(c, "Error getting user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, user)
}

// @Summary		Update user by ID
// @Description	Updates a user's data. Only accessible by admin employees.
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			id path int true "User ID"
// @Param			user body model.UpdateUserRequest true "Updated user information"
// @Success			200 {object} model.ApiResponse "User updated successfully"
// @Failure			400 {object} model.ApiResponse "Invalid request format"
// @Failure			403 {object} model.ApiResponse "Insufficient permissions"
// @Failure			404 {object} model.ApiResponse "User not found"
// @Failure			500 {object} model.ApiResponse "Error updating user"
// @Router			/users/{id} [put]
func UpdateUser(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var req model.UserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request format", http.StatusBadRequest)
		return
	}

	id := c.Param("id")
	var user model.User
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&user).Error; err != nil {
		utils.RespondWithError(c, "User not found", http.StatusNotFound)
		return
	}

	updated := model.User{
		ID:                    user.ID,
		Name:                  req.Name,
		LastName:              req.LastName,
		Phone:                 req.Phone,
		Email:                 req.Email,
		Role:                  req.Role,
		LicenseExpirationDate: req.LicenseExpirationDate,
		IsActive:              true,
		CreatedAt:             user.CreatedAt,
		LastModifiedAt:        time.Now(),
	}

	if err := database.DB.Save(&updated).Error; err != nil {
		utils.RespondWithError(c, "Error updating user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "User updated successfully",
	})
}

// @Summary		    Mark user as inactive
// @Description	    Marks the user as inactive instead of permanently deleting.
// @Tags			users
// @Produce		    json
// @Param			id path int true "User ID"
// @Success		    200 {object} model.ApiResponse "User successfully marked as inactive"
// @Failure		    403 {object} model.ApiResponse "Insufficient permissions"
// @Failure		    500 {object} model.ApiResponse "Error deleting user"
// @Router			/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	id := c.Param("id")

	err := database.DB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error

	if err != nil {
		utils.RespondWithError(c, "Error deleting user", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "User successfully deleted",
	})
}
