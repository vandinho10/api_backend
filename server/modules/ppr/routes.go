// pwd: /app/server/modules/ppr/router.go
package ppr

import (
	"strconv"

	"api/logger"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes adiciona a rota do cálculo de PPR
func RegisterRoutes(router *gin.Engine) {
	group := router.Group("/ppr")
	{
		group.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong PPR"})
		})

		group.GET("/calculate", calculatePPRHandler)
	}
	healthCheck()
}

// calculatePPRHandler lida com a rota de cálculo de PPR
func calculatePPRHandler(c *gin.Context) {
	salaryStr := c.Query("salary")
	pprValueStr := c.Query("ppr_value")
	monthsWorkedValueStr := c.Query("months_worked")

	// Converte os parâmetros de string para float
	salary, err1 := strconv.ParseFloat(salaryStr, 64)
	pprValue, err2 := strconv.ParseFloat(pprValueStr, 64)
	monthsWorked, err3 := strconv.ParseFloat(monthsWorkedValueStr, 64)

	// Verifica se houve erro na conversão ou se o salário e ppr são válidos
	if err1 != nil || err2 != nil || salary <= 0 || pprValue <= 0 {
		c.JSON(400, gin.H{"error": "Parâmetros inválidos"})
		return
	}

	// Verifica o valor de monthsWorked
	if err3 != nil || monthsWorked < 1 || monthsWorked > 12 {
		// Se não for passado ou for inválido, define como 12
		monthsWorked = 12
	}

	// Calcula o PPR
	grossPPR, tax, netPPR := calculatePPR(salary, pprValue, monthsWorked)

	// Prepara a resposta
	response := gin.H{
		"salary":       salary,
		"monthsWorked": monthsWorked,
		"gross_ppr":    grossPPR,
		"tax":          tax,
		"net_ppr":      netPPR,
	}

	// Retorna a resposta
	c.JSON(200, response)
}
func healthCheck() {
	logger.Debug("Módulo PPR está saudável.")
}

func init() {
	logger.Debug("Módulo PPR carregado.")
}
