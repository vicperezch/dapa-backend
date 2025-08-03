package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondWithJSON sends a JSON response with HTTP status 200 (OK).
// The payload parameter can be any data structure to be sent as JSON.
func RespondWithJSON(c *gin.Context, payload any) {
	c.JSON(http.StatusOK, payload)
}

// RespondWithError sends a JSON response with a custom HTTP status code.
// The response contains a success flag set to false and an error message.
func RespondWithError(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"Success": false,
		"Message": message,
	})
}