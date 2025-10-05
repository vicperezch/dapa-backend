package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary      Register a new employee
// @Description  Creates a new employee entry. Only admins are allowed to register employees.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        register body model.RegisterRequest true "Employee registration data"
// @Success      200 {object} model.ApiResponse "Employee registered successfully"
// @Failure      400 {object} model.ApiResponse "Invalid request format"
// @Failure      403 {object} model.ApiResponse "Insufficient permissions"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /users/ [post]
func RegisterHandler(c *gin.Context) {
	var req model.RegisterDTO
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	passwordHash, err := utils.HashString(req.Password)
	if err != nil {
		utils.RespondWithInternalError(c, "Eror registering user")
		return
	}

	user := model.User{
		Name:                  req.Name,
		LastName:              req.LastName,
		Phone:                 req.Phone,
		Email:                 req.Email,
		LicenseExpirationDate: req.LicenseExpirationDate,
		PasswordHash:          passwordHash,
		Role:                  req.Role,
	}

	if err = database.DB.Create(&user).Error; err != nil {
		utils.RespondWithCustomError(
			c,
			http.StatusInternalServerError,
			"Error registering user",
			"Something went wrong",
		)
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, nil, "User registered successfully")
}

// @Summary      Employee login
// @Description  Authenticates an employee and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body model.LoginRequest true "Employee login credentials"
// @Success      200 {object} model.ApiResponse "Login successful, token returned in data"
// @Failure      400 {object} model.ApiResponse "Invalid request format"
// @Failure      401 {object} model.ApiResponse "Invalid credentials"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /login/ [post]
func LoginHandler(c *gin.Context) {
	var req model.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	var user model.User

	err := database.DB.
		Where("email = ? AND is_active = ?", req.Email, true).
		First(&user).Error

	if err != nil || !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.RespondWithCustomError(
			c,
			http.StatusUnauthorized,
			"Email or password is incorrect",
			"Invalid credentials",
		)
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		utils.RespondWithInternalError(c, "Error logging user in")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, token, "User successfully logged in")
}

// @Summary      Reset password request
// @Description  Initiates the process to reset an account's password with a link sent via email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body model.ForgotPasswordDTO true "User email"
// @Success      200 {object} model.ApiResponse "Email sent successfully"
// @Failure      400 {object} model.ApiResponse "Invalid request format"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /auth/forgot [post]
func ForgotPasswordHandler(c *gin.Context) {
	var req model.ForgotPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	var err error
	var user model.User

	err = database.DB.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error sending reset email")
		return
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		utils.RespondWithInternalError(c, "Error sending reset email")
		return
	}

	// Hashear el token y almacenar en la base de datos
	hash, _ := utils.HashString(token)
	expiry := time.Now().Add(30 * time.Minute)
	resetToken := model.ResetToken{
		Token:  hash,
		Expiry: expiry,
		UserID: user.ID,
		IsUsed: false,
	}

	err = database.DB.Create(&resetToken).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error sending reset email")
		return
	}

	link := "http://localhost:5173/reset-password?token=" + token
	emailContent := fmt.Sprintf("<p>Puedes actualizar tu contraseña a través del siguiente link:</p><a href=\"%s\">%s</a>", link, link)

	err = utils.SendEmail(user.Email, "Reestablecimiento de contraseña", emailContent)
	if err != nil {
		utils.RespondWithInternalError(c, "Error sending reset email")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, nil, "Reset email sent")
}

// @Summary      Changes an user's password
// @Description  Uses a reset token to change the account's password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body model.ResetPasswordDTO true "Reset data"
// @Success      200 {object} model.ApiResponse "Password updated successfully"
// @Failure      400 {object} model.ApiResponse "Invalid request format"
// @Failure      401 {object} model.ApiResponse "Token verification failed"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /auth/forgot [post]
func ResetPasswordHandler(c *gin.Context) {
	var req model.ResetPasswordDTO
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	hash, _ := utils.HashString(req.Token)

	var resetToken model.ResetToken
	var user model.User

	err = database.DB.Where("token = ? AND is_used = ?", hash, false).First(&resetToken).Error
	if err != nil {
		utils.RespondWithInternalError(c, "Error resetting password")
		return
	}

	err = database.DB.Where("id = ? AND is_active = ?", resetToken.UserID, true).First(&user).Error
	if err != nil || time.Now().After(resetToken.Expiry) {
		utils.RespondWithInternalError(c, "Error resetting password")
		return
	}

	paswordHash, err := utils.HashString(req.NewPassword)
	if err != nil {
		utils.RespondWithInternalError(c, "Error resetting password")
		return
	}

	user.PasswordHash = paswordHash
	resetToken.IsUsed = true

	database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Save(&user).Error
		if err != nil {
			return err
		}

		err = tx.Save(&resetToken).Error
		if err != nil {
			return err
		}

		return nil
	})

	utils.RespondWithSuccess(c, http.StatusOK, nil, "The user's password has been updated")
}
