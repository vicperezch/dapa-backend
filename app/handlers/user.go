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
func GetUsersHandler(c *gin.Context) {
	var users []model.User

	if err := database.DB.Where("is_active = ?", true).Find(&users).Error; err != nil {
		utils.RespondWithInternalError(c, "Error fetching users")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, users, "Users fetched successfully")
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
func GetUserHandler(c *gin.Context) {
	var user model.User

	id := c.Param("id")
	if err := database.DB.
		Where("id = ? AND is_active = ?", id, true).
		First(&user).Error; err != nil {

		utils.RespondWithInternalError(c, "Error fetching user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user, "User fetched successfully")
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
func UpdateUserHandler(c *gin.Context) {
	var req model.UserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	id := c.Param("id")
	var user model.User
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&user).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusNotFound,
			"User not found",
			"Something went wrong",
		)
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
		utils.RespondWithInternalError(c, "Error updating user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "User updated successfully")
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
func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	err := database.DB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"is_active":  false,
		}).Error

	if err != nil {
		utils.RespondWithInternalError(c, "Error deleting user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "User deleted successfully")
}
