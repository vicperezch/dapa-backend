package utils

import (
	"dapa/app/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RespondWithJSON sends a JSON response with HTTP status 200 (OK).
// The payload parameter can be any data structure to be sent as JSON.
func RespondWithSuccess(c *gin.Context, status int, data any, message string) {
	c.JSON(status, model.ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// RespondWithError sends a JSON response with a custom HTTP status code.
// The response contains a success flag set to false and an error message.
func RespondWithError(c *gin.Context, status int, err error, message string) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		errorMessages := make([]string, len(ve))

		for i, fe := range ve {
			errorMessages[i] = getTagMessage(fe.Tag())
		}

		c.JSON(status, model.ApiResponse{
			Success: false,
			Message: message,
			Data:    nil,
			Errors:  errorMessages,
		})
		return
	}

	c.JSON(status, model.ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  []string{"An unknown error ocurred"},
	})
}

func RespondWithCustomError(c *gin.Context, status int, err string, message string) {
	c.JSON(status, model.ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  []string{err},
	})
}

func RespondWithInternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, model.ApiResponse{
		Success: false,
		Message: "Something went wrong",
		Data:    nil,
		Errors:  []string{err},
	})
}

func RespondWithUnathorizedError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, model.ApiResponse{
		Success: false,
		Message: "Only administrators can perform this action",
		Data:    nil,
		Errors:  []string{"Insufficient permissions"},
	})
}

