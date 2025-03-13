// // pwd: /app/server/modules/login/middleware/traefik_auth.go
package middleware

// import (
// 	// "errors"
// 	"net/http"
// 	"strings"

// 	"api/logger"
// 	"api/server/modules/login/services"
// 	"github.com/gin-gonic/gin"
// )

// // TraefikAuthMiddleware verifica se a requisição contém um token JWT válido antes de liberar o acesso
// func TraefikAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Obtém o token JWT do cabeçalho Authorization
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			logger.Warn("Tentativa de acesso sem token JWT")
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT ausente"})
// 			c.Abort()
// 			return
// 		}

// 		// Extrai o token removendo "Bearer "
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == "" {
// 			logger.Warn("Token JWT inválido no cabeçalho Authorization")
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
// 			c.Abort()
// 			return
// 		}

// 		// Valida o token JWT
// 		userID, err := services.ValidateJWT(tokenString)
// 		if err != nil {
// 			logger.Warn("Falha na autenticação JWT: %s", err.Error())
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
// 			c.Abort()
// 			return
// 		}

// 		// Define o usuário autenticado no contexto
// 		c.Set("user_id", userID)

// 		// Continua para o serviço protegido
// 		c.Next()
// 	}
// }
