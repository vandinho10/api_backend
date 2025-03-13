// pwd: /app/server/modules/finance/router.go
package finance

import (
	"api/logger"
	"api/utils"

	"github.com/gin-gonic/gin"
)

var BearerToken = utils.GetEnv("BEARER_PROTECTED_PATHS")
var FinancePath = utils.GetEnv("FINANCE_PATH")
var FinanceCsv = utils.GetEnv("FINANCE_CSV")
var FinanceCsvDb = utils.GetEnv("FINANCE_CSV_DB")

// RegisterRoutes adiciona a rota do cálculo de PPR
func RegisterRoutes(router *gin.Engine) {
	group := router.Group(FinancePath)
	protected := group.Group("/")
	protected.Use(AuthMiddleware()) // Protege as rotas subsequentes
	{
		protected.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		protected.POST(FinanceCsv, handleExtract)
		protected.POST(FinanceCsvDb, handleExtractDB)
	}
}

func init() {
	logger.Debug("Módulo Finance carregado.")
}
