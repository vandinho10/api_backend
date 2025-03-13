// pwd: /app/server/modules/ppr/service.go
package ppr

import (
	"math"

	"api/logger"
)

// calculatePPR calcula o valor do PPR e o imposto retido corretamente com a tabela progressiva
func calculatePPR(salary, pprValue float64, monthsWorked float64) (grossPPR, tax, netPPR float64) {
	if salary <= 0 || pprValue <= 0 {
		logger.Warn("Valores inválidos fornecidos para cálculo do PPR - Salary: %.2f, PPR Value: %.2f", salary, pprValue)
		return 0, 0, 0
	}

	// Calcula o PPR bruto mensal
	grossPPRMonthly := salary * pprValue / 12

	// Ajusta o PPR bruto de acordo com os meses trabalhados
	grossPPR = grossPPRMonthly * monthsWorked

	// Calcula o imposto sobre o valor ajustado do PPR
	tax = calculateTax(grossPPR)

	// Calcula o PPR líquido
	netPPR = grossPPR - tax

	logger.Debug("Cálculo do PPR realizado - Gross PPR: %.2f, Tax: %.2f, Net PPR: %.2f", grossPPR, tax, netPPR)

	return
}

// calculateTax aplica a tabela progressiva do IR corretamente
func calculateTax(amount float64) float64 {
	var tax float64

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

	for _, bracket := range taxBrackets {
		if amount > bracket.limit {
			tax += (bracket.limit - previousLimit) * bracket.rate
		} else {
			tax += (amount - previousLimit) * bracket.rate
			break
		}
		previousLimit = bracket.limit
	}

	logger.Debug("Cálculo do imposto realizado - Amount: %.2f, Tax: %.2f", amount, tax)

	return tax
}
