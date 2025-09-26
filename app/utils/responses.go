package utils

import (
	"dapa/app/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Envía una respuesta en JSON en caso de éxito
// Recibe el código HTTP, la información y el mensaje a enviar
func RespondWithSuccess(c *gin.Context, status int, data any, message string) {
	c.JSON(status, model.ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// Envía una respuesta en JSON en caso de fallo
// Recibe el código HTTP, el error encontrado y el mensaje como parámetro
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

// Envía una respuesta personalizada en JSON en caso de fallo
// Recibe el código HTTP, el error a mostrar y el mensaje como parámetro
func RespondWithCustomError(c *gin.Context, status int, err string, message string) {
	c.JSON(status, model.ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  []string{err},
	})
}

// Envía una respuesta en JSON en caso de fallo interno del programa
// Recibe el error a mostrar como parámetro
func RespondWithInternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, model.ApiResponse{
		Success: false,
		Message: "Something went wrong",
		Data:    nil,
		Errors:  []string{err},
	})
}

// Envía una respuesta personalizada en JSON en caso de usuario no autorizado
func RespondWithUnathorizedError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, model.ApiResponse{
		Success: false,
		Message: "Only administrators can perform this action",
		Data:    nil,
		Errors:  []string{"Insufficient permissions"},
	})
}
