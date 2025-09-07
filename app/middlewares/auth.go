package middlewares

import (
	"net/http"

	"dapa/app/model"
	"dapa/app/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware that checks if a valid JWT token is present.
// It extracts the token from the Authorization header, validates it,
// and stores the claims in the Gin context for further use.
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

// RoleRequired is a Gin middleware that verifies if the user has one of the allowed roles.
// It must be used after AuthMiddleware, since it relies on the token claims being present in the context.
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

// extractToken extracts the JWT token from the Authorization header.
// It expects the format: "Bearer <token>"
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}

