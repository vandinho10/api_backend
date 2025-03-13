// pwd: /app/server/healthCheck.go
package server

import (
	"api/db"
	"api/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Função que realiza o Health Check
func HealthCheck(c *gin.Context) {
	// Simulação de verificação de dependências, como banco de dados
	dbStatus := checkDatabaseConnection()

	if dbStatus {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "unhealthy",
			"message": "Database connection failed",
		})
	}
}

// Função para verificar a conexão com o banco de dados
func checkDatabaseConnection() bool {
	// Tenta estabelecer a conexão com o banco de dados
	dbConnection, err := db.DbConnection()

	// Se a conexão falhar, retorna falso
	if err != nil {
		logger.Error("Database connection failed: %v", err) // Loga o erro de conexão
		return false
	}

	// Se a conexão for bem-sucedida, fecha a conexão antes de retornar
	defer dbConnection.Close()

	// Se não houver erro, a conexão foi bem-sucedida
	return true
}
