// pwd: /app/server/modules/login/services/auth_service.go
package services

import (
	"errors"

	"api/logger"
	"api/server/modules/login/auth_utils"
	"api/server/modules/login/models"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest representa a estrutura da requisição de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest representa a estrutura da requisição de registro
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// Login autentica o usuário e gera um token JWT
func Login(req LoginRequest) (string, error) {
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		logger.Warn("Usuário não encontrado: %v", req.Email)
		return "", errors.New("credenciais inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("Senha incorreta para: %v", req.Email)
		return "", errors.New("credenciais inválidas")
	}

	// Aqui você pode usar a função GenerateJWT de jwt_utils.go
	token, err := auth_utils.GenerateJWT(user.ID, user.Username, user.AccessLevel) // Passa apenas o ID, pois os outros dados estão definidos em Claims
	if err != nil {
		logger.Error("Erro ao gerar token JWT: %v", err)
		return "", errors.New("erro ao gerar token")
	}

	return token, nil
}

// Register cria um novo usuário no banco de dados
func Register(req RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Erro ao criptografar senha: %v", err)
		return errors.New("erro ao registrar usuário")
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := models.CreateNewUser(user); err != nil {
		logger.Error("Erro ao salvar usuário no banco: %v", err)
		return errors.New("erro ao registrar usuário")
	}

	return nil
}
