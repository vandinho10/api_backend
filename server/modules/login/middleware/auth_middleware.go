// pwd: /app/server/modules/login/middleware/auth_middleware.go
package middleware

import (
	"net/http"
	"strings"

	"api/logger"
	"api/server/modules/login/auth_utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica e valida o token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// O token vem no formato: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do token inválido"})
			c.Abort()
			return
		}

		// Valida o token e extrai os claims
		claims, err := auth_utils.ValidateAndExtractClaims(tokenParts[1])
		if err != nil {
			logger.Warn("Falha ao validar token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Define os valores extraídos do token no contexto
		c.Set("user_id", claims.ID)
		// c.Set("username", claims.Username)
		c.Set("access_level", claims.AccessLevel)

		c.Next()
	}
}
