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

// AuthenticateUser autentica um usuário com base nos dados fornecidos na requisição JSON e retorna um token JWT.
//
// Parâmetros:
// - c: *gin.Context - Contexto da requisição.
//
// Respostas:
// - 200 OK: Retorna o token JWT, informações do usuário e tempo restante até a expiração.
// - 400 Bad Request: Se os dados da requisição estiverem inválidos.
// - 401 Unauthorized: Se as credenciais forem inválidas.
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

// IsLoggedIn verifica se o usuário está autenticado e retorna o tempo restante de expiração do token.
//
// Parâmetros:
// - c: *gin.Context - Contexto da requisição.
//
// Respostas:
// - 200 OK: Retorna se o usuário está logado e o tempo restante do token.
// - 401 Unauthorized: Se o token não for fornecido, for inválido ou estiver expirado.
func IsLoggedIn(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	if models.IsTokenBlacklisted(tokenString) {
		logger.Warn("Tentativa de uso de token inválido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	timeRemaining, err := auth_utils.CalculateTokenExpirationTime(tokenString)
	logger.Debug("Tempo restante até a expiração do token: %v", timeRemaining)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Erro ao calcular o tempo de expiração do token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logged_in":      true,
		"time_remaining": timeRemaining.String(),
	})
}

// LogoutUser realiza o logout do usuário, invalidando o token JWT.
//
// Parâmetros:
// - c: *gin.Context - Contexto da requisição.
//
// Respostas:
// - 200 OK: Se o logout for bem-sucedido.
// - 400 Bad Request: Se o token não for fornecido ou houver erro ao processar o logout.
// - 401 Unauthorized: Se o token for inválido ou já estiver na blacklist.
func LogoutUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	if models.IsTokenBlacklisted(tokenString) {
		logger.Warn("Tentativa de uso de token inválido")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	expirationDuration, err := auth_utils.CalculateTokenExpirationTime(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao calcular expiração do token"})
		return
	}

	expirationTime := time.Now().Add(expirationDuration)

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
//
// Parâmetros:
// - c: *gin.Context - Contexto da requisição.
//
// Respostas:
// - 201 Created: Se o usuário for criado com sucesso.
// - 400 Bad Request: Se os dados forem inválidos ou o usuário já existir.
// - 500 Internal Server Error: Se ocorrer um erro ao criar o usuário.
func AddNewUser(c *gin.Context) {
	logger.Debug("Adicionando novo usuário")

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		logger.Warn("Token JWT não fornecido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token JWT não fornecido"})
		return
	}

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		logger.Error("Erro ao validar dados do usuário: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	err := models.CheckUserExists(newUser.Username)
	if err != nil {
		logger.Warn("Usuário já existe")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuário já existe"})
		return
	}

	err = models.CreateNewUser(newUser)
	if err != nil {
		logger.Error("Erro ao criar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	logger.Info("Usuário registrado com sucesso")
	c.JSON(http.StatusCreated, gin.H{"message": "Usuário registrado com sucesso"})
}
