// pwd: /app/server/modules/finance/router.go
package finance

import (
	"api/logger"
	"api/utils"

	"github.com/gin-gonic/gin"
)

// Variáveis de ambiente carregadas
var BearerToken = utils.GetEnv("BEARER_PROTECTED_PATHS")
var FinancePath = utils.GetEnv("FINANCE_PATH")
var FinanceCsv = utils.GetEnv("FINANCE_CSV")
var FinanceCsvDb = utils.GetEnv("FINANCE_CSV_DB")

// RegisterRoutes registra as rotas do módulo Finance.
//
// Esta função adiciona as rotas para o módulo de finanças, incluindo rotas protegidas
// que requerem autenticação via middleware. A rota de /ping serve para verificar se
// o módulo de finanças está respondendo corretamente. As rotas de upload de CSV e extração de dados
// também são configuradas aqui.
//
// Parâmetros:
//   - router (*gin.Engine): A instância do router Gin, usada para adicionar rotas.
func RegisterRoutes(router *gin.Engine) {
	// Cria um grupo de rotas com o caminho configurado na variável FinancePath
	group := router.Group(FinancePath)

	// Cria um grupo protegido com middleware de autenticação
	protected := group.Group("/")
	protected.Use(AuthMiddleware()) // Aplica o middleware de autenticação

	// Adiciona as rotas protegidas
	{
		// Rota de verificação
		protected.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - Finance"})
		})

		// Rota para upload de arquivos CSV para extração de dados financeiros
		protected.POST(FinanceCsv, handleExtract)

		// Rota para extração de dados financeiros diretamente do banco de dados
		protected.POST(FinanceCsvDb, handleExtractDB)
	}
}

// init inicializa o módulo Finance.
//
// Esta função é chamada automaticamente ao carregar o pacote. Ela apenas registra um log
// indicando que o módulo de finanças foi carregado com sucesso.
func init() {
	logger.Debug("Módulo Finance carregado.")
}
