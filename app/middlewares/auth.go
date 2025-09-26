package middlewares

import (
	"net/http"

	"dapa/app/model"
	"dapa/app/utils"

	"github.com/gin-gonic/gin"
)

// Middleware para determinar si un token JWT es válido
// Si es válido, almacena las claims para uso posterior
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		if tokenString == "" {
			utils.RespondWithUnathorizedError(c)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithUnathorizedError(c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// Middleware para verificar si el usuario posee los roles requeridos
// Recibe los roles requeridos como parámetro
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			utils.RespondWithUnathorizedError(c)
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*model.EmployeeClaims)
		if !ok {
			utils.RespondWithCustomError(c, http.StatusBadRequest, "Invalid claims format", "Invalid request format")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.RespondWithUnathorizedError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Extrae el token JWT del authorization header
// Retorna el token como string
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}
