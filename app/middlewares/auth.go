package middlewares

import (
	"net/http"
	
	"dapa/app/utils"
	"dapa/app/model"

	"github.com/gin-gonic/gin"
)

//Middleware de autenticación que corrobora que el token exista y sea válido
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		if tokenString == "" {
			utils.RespondWithError(c, "Authorization token required", http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(c, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Almacenar claims en el contexto de Gin
		c.Set("claims", claims)
		c.Next()
	}
}

//Middleware que verifica roles
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			utils.RespondWithError(c, "Claims not found", http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*model.EmployeeClaims)
		if !ok {
			utils.RespondWithError(c, "Invalid claims format", http.StatusInternalServerError)
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
			utils.RespondWithError(c, "Insufficient permissions", http.StatusForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

//Función que obtiene el token del header y retira los caracteres innecesarios
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}