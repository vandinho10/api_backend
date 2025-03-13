// pwd: /app/server/modules/ppr/router.go
package ppr

import (
	"api/logger"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra as rotas do módulo de PPR (Plano de Participação nos Resultados).
//
// Esta função adiciona as rotas necessárias ao grupo "/ppr", incluindo o endpoint "/ping"
// para verificação de funcionamento e o endpoint "/calculate" para o cálculo do PPR.
//
// Parâmetros:
//   - router (*gin.Engine): A instância do roteador Gin onde as rotas serão registradas.
func RegisterRoutes(router *gin.Engine) {
	// Cria um grupo de rotas para o módulo PPR
	group := router.Group("/ppr")
	{
		// Rota de verificação de funcionamento (ping)
		group.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong PPR"})
		})

		// Rota para calcular o PPR
		group.GET("/calculate", calculatePPRHandler)
	}

	// Chama a função de verificação de saúde do módulo
	healthCheck()
}

// healthCheck loga que o módulo PPR está saudável e funcionando corretamente.
func healthCheck() {
	// Loga uma mensagem de debug indicando que o módulo está funcionando
	logger.Debug("Módulo PPR está saudável.")
}

// init é uma função especial que é chamada automaticamente quando o pacote é carregado.
// Aqui, usamos para logar que o módulo PPR foi carregado com sucesso.
func init() {
	// Loga uma mensagem de debug indicando que o módulo foi carregado
	logger.Debug("Módulo PPR carregado.")
}
