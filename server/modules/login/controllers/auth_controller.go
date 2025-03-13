// pwd: /app/server/modules/login/controllers/auth_controller.go
package controllers

import (
	"net/http"
	"time"

	"api/logger"
	"api/server/modules/login/auth_utils"
	"api/server/modules/login/models"

	"github.com/gin-gonic/gin"
)

// AuthenticateUser autentica um usuário com base nos dados fornecidos na requisição JSON.
func AuthenticateUser(c *gin.Context) {
	var loginData models.LoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
		return
	}

	authResponse, err := models.Authenticate(loginData)
	if err != nil {
		logger.Warn("Tentativa de login falhou para o usuário %s", loginData.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	// Retorna a resposta completa com token, usuário e tempo restante
	c.JSON(http.StatusOK, gin.H{
		"token": authResponse.Token,
		"user": gin.H{
			"id":           authResponse.User.ID,
			"username":     authResponse.User.Username,
			"email":        authResponse.User.Email,
			"access_level": authResponse.User.AccessLevel,
		},
		"time_remaining": authResponse.TimeRemaining,
	})
}

// IsLoggedIn verifica se o usuário está autenticado e retorna informações sobre o usuário e o tempo restante até a expiração do token JWT.
func IsLoggedIn(c *gin.Context) {
	// Recupera o token JWT do header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	// Recupera os dados do usuário do contexto
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Verifica se o token está na blacklist
	if models.IsTokenBlacklisted(tokenString) {
		logger.Warn("Tentativa de uso de token inválido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	// Calcula o tempo restante até a expiração do token
	timeRemaining, err := auth_utils.CalculateTokenExpirationTime(tokenString)
	logger.Debug("Tempo restante até a expiração do token: %v", timeRemaining)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro ao calcular o tempo de expiração do token"})
		return
	}

	// Retorna a resposta no formato desejado
	c.JSON(http.StatusOK, gin.H{
		"logged_in":      true,
		"time_remaining": timeRemaining.String(), // Inclui o tempo restante até a expiração do token
	})
}

// LogoutUser realiza o logout do usuário, removendo o token JWT.
func LogoutUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	// Verifica se o token está na blacklist
	if models.IsTokenBlacklisted(tokenString) {
		logger.Warn("Tentativa de uso de token inválido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	// Calcula o tempo de expiração do token
	expirationDuration, err := auth_utils.CalculateTokenExpirationTime(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao calcular expiração do token"})
		return
	}

	// Converte para um time.Time válido
	expirationTime := time.Now().Add(expirationDuration)

	// Adiciona o token à blacklist
	err = models.AddTokenToBlacklist(tokenString, expirationTime)
	if err != nil {
		logger.Error("Erro ao adicionar token à blacklist: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar logout"})
		return
	}

	logger.Info("Usuário deslogado com sucesso")
	c.JSON(http.StatusOK, gin.H{"message": "Logout realizado. Remova o token do cliente."})
}

// AddNewUser cria um novo usuário no sistema com base nos dados fornecidos na requisição JSON.
func AddNewUser(c *gin.Context) {
	logger.Debug("Adicionando novo usuário")

	logger.Debug("Validando token JWT")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logger.Warn("Token JWT não fornecido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	logger.Debug("Validando dados do usuário")
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		logger.Error("Erro ao validar dados do usuário: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	logger.Debug("Verificando se o usuário já existe")
	err := models.CheckUserExists(newUser.Username)
	if err != nil {
		logger.Warn("Usuário já existe")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário já existe"})
		return
	}

	logger.Debug("Criando novo usuário")
	err = models.CreateNewUser(newUser)
	if err != nil {
		logger.Error("Erro ao criar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	logger.Info("Usuário registrado com sucesso")
	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}

// // ValidateTokenHandler valida o token JWT enviado
// func ValidateTokenHandler(c *gin.Context) {
// 	userID, _ := c.Get("user_id")
// 	c.JSON(http.StatusOK, gin.H{"message": "Token válido", "user_id": userID})
// }
