package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithJSON(c *gin.Context, payload any)  {
  c.JSON(http.StatusOK, payload)
}

func RespondWithError(c *gin.Context, message string, status int)  {
  c.JSON(status, gin.H{
    "Success": false,
    "Message": message,
  })
}