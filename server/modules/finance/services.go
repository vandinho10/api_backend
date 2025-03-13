// pwd: /app/server/modules/finance/router.go
package finance

import (
	"api/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware é um middleware que verifica se o token Bearer fornecido é válido.
//
// Este middleware verifica o cabeçalho "Authorization" da requisição para garantir que
// o token Bearer seja fornecido e seja válido. Caso o token esteja ausente ou seja inválido,
// ele retorna um erro com o código de status adequado (401 ou 403) e interrompe a requisição.
//
// Retorna:
//   - gin.HandlerFunc: A função de middleware para ser usada no Gin.
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
//
// Esta função configura os cabeçalhos de resposta para o arquivo CSV e o envia
// para o cliente. O arquivo será enviado com o tipo MIME correto (text/csv) e
// será exibido inline no navegador.
//
// Parâmetros:
//   - c (*gin.Context): O contexto da requisição do Gin.
//   - filePath (string): O caminho do arquivo a ser enviado.
func fileExtract(c *gin.Context, filePath string) {
	// Define os cabeçalhos da resposta
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "inline")
	// Envia o arquivo especificado para o cliente
	c.File(filePath)
}

// handleExtract lida com a requisição para servir o arquivo extract.csv.
//
// Esta função é acionada quando o cliente solicita o arquivo extract.csv. Ela
// chama a função fileExtract para enviar o arquivo correto.
//
// Parâmetros:
//   - c (*gin.Context): O contexto da requisição do Gin.
func handleExtract(c *gin.Context) {
	fileExtract(c, "/app/server/modules/finance/files/files/extract.csv")
}

// handleExtractDB lida com a requisição para servir o arquivo extract_db.csv.
//
// Esta função é acionada quando o cliente solicita o arquivo extract_db.csv. Ela
// chama a função fileExtract para enviar o arquivo correto.
//
// Parâmetros:
//   - c (*gin.Context): O contexto da requisição do Gin.
func handleExtractDB(c *gin.Context) {
	fileExtract(c, "/app/server/modules/finance/files/files/extract_db.csv")
}

// init inicializa o módulo Finance.
//
// Esta função é chamada automaticamente ao carregar o pacote. Ela apenas registra um log
// indicando que o módulo de finanças foi carregado com sucesso.
func init() {
	logger.Debug("Módulo Finance carregado.")
}
