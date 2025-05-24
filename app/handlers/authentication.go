package handlers

import (
	"dapa/app/model"
	"dapa/app/utils"
	"dapa/database"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary		Register a new employee.
// @Description	Creates a new employee entry in the database.
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		register body model.RegisterRequest true "Required information to register user."
// @Success		200	{object} model.ApiResponse "Employee registered successfully."
// @Failure		400	{object} model.ApiResponse "Request did not pass the validations."
// @Failure		500	{object} model.ApiResponse "Error when trying to register user."
// @Router		/users/ [post]
func RegisterHandler(c *gin.Context) {
	claims := c.MustGet("claims").(*model.EmployeeClaims)
	
	if claims.Role != "admin" {
		utils.RespondWithError(c,"Insufficient permissions",http.StatusForbidden )
		return
	}

	var req model.RegisterRequest
	var err error
	
	if err = c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request: ", err)

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
		log.Println("Error hashing password: ", err)
		utils.RespondWithError(c, "Error hashing password", http.StatusInternalServerError)
		return
	}

	employee := model.Employee{
		User: model.User{
			Name:     req.Name,
			LastName: req.LastName,
			Phone:    req.Phone,
			Email:    req.Email,
		},
		Password: passwordHash,
		Role:     req.Role,
	}

	if err = database.DB.Create(&employee).Error; err != nil {
		utils.RespondWithError(c, "Error registrando al usuario", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "User created successfully",
	})
}

// @Summary		Login for employees.
// @Description	Authenticates an employee and returns a JWT token.
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		login body model.LoginRequest true "Credentials for login"
// @Success		200	{object} model.ApiResponse "Login successful"
// @Failure		400	{object} model.ApiResponse "Invalid request format"
// @Failure		401	{object} model.ApiResponse "Invalid email or password"
// @Router		/login/ [post]
func LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error parsing request: ", err)
		utils.RespondWithError(c, "Invalid request", http.StatusBadRequest)
		return
	}

	var employee model.Employee

	err := database.DB.Preload("User").
		Joins("JOIN users ON users.id = employees.user_id").
		Where("users.email = ?", req.Email).
		First(&employee).Error

	if err != nil {
		log.Println("Error finding user: ", err)
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(req.Password, employee.Password) {
		utils.RespondWithError(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(&employee)
	if err != nil {
		utils.RespondWithError(c,"Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(c, model.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data: token,
	})
}
