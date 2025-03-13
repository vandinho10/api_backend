// pwd: /app/server/modules/ppr/service.go
package ppr

import (
	"math"
	"strconv"

	"api/logger"

	"github.com/gin-gonic/gin"
)

// calculatePPRHandler lida com a requisição para calcular o PPR (Plano de Participação nos Resultados).
//
// A função recebe os parâmetros "salary", "ppr_value" e "months_worked" da query string,
// realiza a validação e conversão dos valores, e então chama a função de cálculo do PPR.
// No final, retorna o resultado do cálculo em formato JSON.
//
// Parâmetros:
//   - c (*gin.Context): O contexto da requisição, utilizado para acessar os parâmetros e enviar a resposta.
func calculatePPRHandler(c *gin.Context) {
	// Obtém os parâmetros passados na query string
	salaryStr := c.Query("salary")
	pprValueStr := c.Query("ppr_value")
	monthsWorkedValueStr := c.Query("months_worked")

	// Converte os parâmetros de string para float
	salary, err1 := strconv.ParseFloat(salaryStr, 64)
	pprValue, err2 := strconv.ParseFloat(pprValueStr, 64)
	monthsWorked, err3 := strconv.ParseFloat(monthsWorkedValueStr, 64)

	// Verifica se houve erro na conversão ou se o salário e PPR são inválidos
	if err1 != nil || err2 != nil || salary <= 0 || pprValue <= 0 {
		c.JSON(400, gin.H{"error": "Parâmetros inválidos"})
		return
	}

	// Verifica o valor de monthsWorked e define um valor padrão de 12 se inválido
	if err3 != nil || monthsWorked < 1 || monthsWorked > 12 {
		monthsWorked = 12
	}

	// Calcula o PPR, incluindo imposto e PPR líquido
	grossPPR, tax, netPPR := calculatePPR(salary, pprValue, monthsWorked)

	// Prepara a resposta com os resultados do cálculo
	response := gin.H{
		"salary":       salary,
		"monthsWorked": monthsWorked,
		"gross_ppr":    grossPPR,
		"tax":          tax,
		"net_ppr":      netPPR,
	}

	// Retorna a resposta em formato JSON
	c.JSON(200, response)
}

// calculatePPR calcula o valor bruto do PPR, o imposto sobre o valor e o PPR líquido.
//
// A função utiliza o salário, o valor do PPR e os meses trabalhados para calcular o valor
// bruto do PPR. Em seguida, aplica a função de cálculo de imposto e retorna o valor líquido.
//
// Parâmetros:
//   - salary (float64): O salário do funcionário.
//   - pprValue (float64): O valor do PPR acordado.
//   - monthsWorked (float64): O número de meses trabalhados no período.
//
// Retorna:
//   - grossPPR (float64): O valor bruto do PPR calculado.
//   - tax (float64): O valor do imposto calculado.
//   - netPPR (float64): O valor líquido do PPR após a dedução do imposto.
func calculatePPR(salary, pprValue float64, monthsWorked float64) (grossPPR, tax, netPPR float64) {
	// Verifica se os valores de salário e PPR são válidos
	if salary <= 0 || pprValue <= 0 {
		logger.Warn("Valores inválidos fornecidos para cálculo do PPR - Salary: %.2f, PPR Value: %.2f", salary, pprValue)
		return 0, 0, 0
	}

	// Calcula o PPR bruto mensal
	grossPPRMonthly := salary * pprValue / 12

	// Ajusta o PPR bruto com base nos meses trabalhados
	grossPPR = grossPPRMonthly * monthsWorked

	// Calcula o imposto sobre o valor do PPR
	tax = calculateTax(grossPPR)

	// Calcula o valor líquido do PPR
	netPPR = grossPPR - tax

	// Loga o resultado do cálculo do PPR
	logger.Debug("Cálculo do PPR realizado - Gross PPR: %.2f, Tax: %.2f, Net PPR: %.2f", grossPPR, tax, netPPR)

	return
}

// calculateTax calcula o imposto a ser retido sobre o valor do PPR, aplicando a tabela progressiva do IR.
//
// A função percorre as faixas da tabela de impostos e calcula o imposto de acordo com o valor do PPR.
// A tabela de impostos é composta por limites de faixa, taxa de imposto e dedução para cada faixa.
//
// Parâmetros:
//   - amount (float64): O valor do PPR que será utilizado no cálculo do imposto.
//
// Retorna:
//   - tax (float64): O valor do imposto calculado.
func calculateTax(amount float64) float64 {
	var tax float64

	// Definição das faixas de imposto com limites, taxas e deduções
	taxBrackets := []struct {
		limit     float64
		rate      float64
		deduction float64
	}{
		{7640.80, 0.00, 0.00},
		{9922.28, 0.075, 573.06},
		{13167.00, 0.15, 1317.23},
		{16380.38, 0.225, 2304.76},
		{math.MaxFloat64, 0.275, 3123.78},
	}

	previousLimit := 0.0

	// Aplica a tabela progressiva de imposto
	for _, bracket := range taxBrackets {
		if amount > bracket.limit {
			tax += (bracket.limit - previousLimit) * bracket.rate
		} else {
			tax += (amount - previousLimit) * bracket.rate
			break
		}
		previousLimit = bracket.limit
	}

	// Loga o valor do imposto calculado
	logger.Debug("Cálculo do imposto realizado - Amount: %.2f, Tax: %.2f", amount, tax)

	return tax
}
