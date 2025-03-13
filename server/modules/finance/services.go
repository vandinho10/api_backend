// pwd: /app/server/modules/finance/router.go
package finance

import (
	"api/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware para verificar o token Bearer
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pega o cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")

		// Verifica se o cabeçalho está presente e se é um Bearer token
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Autorização necessária"})
			c.Abort()
			return
		}

		// Extrai o token do cabeçalho
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Verifica se o token é válido
		if token != BearerToken {
			c.JSON(403, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Se o token for válido, segue com a requisição
		c.Next()
	}
}

// fileExtract envia o arquivo especificado para o cliente com os headers apropriados.
func fileExtract(c *gin.Context, filePath string) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "inline")
	c.File(filePath)
}

// handleExtract manipula a requisição para servir o arquivo extract.csv.
func handleExtract(c *gin.Context) {
	fileExtract(c, "/app/server/modules/finance/files/files/extract.csv")
}

// handleExtractDB manipula a requisição para servir o arquivo extract_db.csv.
func handleExtractDB(c *gin.Context) {
	fileExtract(c, "/app/server/modules/finance/files/files/extract_db.csv")
}

func init() {
	logger.Debug("Módulo Finance carregado.")
}
