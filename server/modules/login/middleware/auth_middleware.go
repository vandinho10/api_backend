// pwd: /app/server/modules/login/middleware/auth_middleware.go

package middleware

import (
	"net/http"
	"strings"

	"api/logger"
	"api/server/modules/login/auth_utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware é um middleware de autenticação que verifica e valida o token JWT presente no cabeçalho da requisição.
//
// Funcionamento:
// - O token deve ser fornecido no cabeçalho "Authorization" no formato: "Bearer <token>".
// - Caso o token não seja fornecido ou esteja em um formato inválido, a requisição é abortada com status 401 (Unauthorized).
// - Se o token for válido, os dados do usuário são extraídos e armazenados no contexto da requisição.
//
// Uso:
// router.Use(AuthMiddleware())
//
// Respostas:
// - 401 Unauthorized: Se o token não for fornecido, estiver inválido ou expirado.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		// O token deve seguir o formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato do token inválido"})
			c.Abort()
			return
		}

		// Valida o token e extrai os claims (dados do usuário)
		claims, err := auth_utils.ValidateAndExtractClaims(tokenParts[1])
		if err != nil {
			logger.Warn("Falha ao validar token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Armazena os dados do usuário extraídos do token no contexto da requisição
		c.Set("user_id", claims.ID)
		// c.Set("username", claims.Username) // Descomentar se necessário
		c.Set("access_level", claims.AccessLevel)

		// Prossegue para a próxima etapa da requisição
		c.Next()
	}
}
