package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary      Register a new employee
// @Description  Creates a new employee entry. Only admins are allowed to register employees.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        register body model.RegisterRequest true "Employee registration data"
// @Success      200 {object} model.ApiResponse "Employee registered successfully"
// @Failure      400 {object} model.ApiResponse "Invalid input data"
// @Failure      403 {object} model.ApiResponse "Insufficient permissions"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /users/ [post]
func RegisterHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)

	if claims.Role != "admin" {
		utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
		return
	}

	var req model.RegisterDTO
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
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
		utils.RespondWithError(c, "Error hashing password", http.StatusInternalServerError)
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
		utils.RespondWithError(c, "Error registrando al usuario", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "User created successfully",
	})
}

// @Summary      Employee login
// @Description  Authenticates an employee and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body model.LoginRequest true "Employee login credentials"
// @Success      200 {object} model.ApiResponse "Login successful, token returned in data"
// @Failure      400 {object} model.ApiResponse "Invalid request format"
// @Failure      401 {object} model.ApiResponse "Invalid email or password"
// @Failure      500 {object} model.ApiResponse "Internal server error"
// @Router       /login/ [post]
func LoginHandler(c *gin.Context) {
	var req model.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request", http.StatusBadRequest)
		return
	}

	var user model.User

	err := database.DB.
		Where("email = ? AND is_active = ?", req.Email, true).
		First(&user).Error

	if err != nil {
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		utils.RespondWithError(c, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    token,
	})
}

func ResetLinkHandler(c *gin.Context) {
	var req model.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request", http.StatusBadRequest)
		return
	}

	// token := utils.GenerateResetToken(req.Email, getPasswordHash(req.Email))
	// TODO: send token via email
}

func PasswordResetHandler(c *gin.Context) {
	var req model.NewPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, "Invalid request", http.StatusBadRequest)
		return
	}

	login, err := utils.VerifyResetToken(req.Token, getPasswordHash)
	if err != nil {
		utils.RespondWithError(c, "Token verification failed", http.StatusUnauthorized)
		return
	}

	paswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.RespondWithError(c, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var user model.User
	database.DB.Where("is_active = ? and email = ?", true, login).Scan(&user)

	user.PasswordHash = paswordHash
	database.DB.Save(&user)

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Password updated succesfully",
	})
}

func getPasswordHash(email string) ([]byte, error) {
	var hash []byte
	database.DB.
		Table("employees").
		Select("employees.password").
		Joins("left join users on users.id = employees.user_id").
		Where("is_active = ? and users.email = ?", true, email).
		Scan(&hash)

	return hash, nil
}
