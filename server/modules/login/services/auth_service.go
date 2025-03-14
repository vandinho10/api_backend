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
	Email    string `json:"email" binding:"required,email"` // O email é obrigatório e deve ter formato válido
	Password string `json:"password" binding:"required"`    // A senha é obrigatória
}

// RegisterRequest representa a estrutura da requisição de registro
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`    // O email é obrigatório e deve ter formato válido
	Password string `json:"password" binding:"required,min=6"` // A senha é obrigatória e deve ter no mínimo 6 caracteres
	Name     string `json:"name" binding:"required"`           // O nome do usuário é obrigatório
}

// Login autentica o usuário e gera um token JWT
//
// Parâmetros:
// - req: A requisição de login contendo email e senha do usuário.
//
// Retorno:
// - string: O token JWT gerado para o usuário.
// - error: Retorna erro caso as credenciais estejam incorretas ou se houver falha na geração do token.
//
// Detalhes:
// - Verifica se o usuário existe e se a senha fornecida corresponde ao hash armazenado.
// - Se as credenciais forem válidas, um token JWT é gerado e retornado.
func Login(req LoginRequest) (string, error) {
	// Recupera o usuário pelo email
	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		logger.Warn("Usuário não encontrado: %v", req.Email)
		return "", errors.New("credenciais inválidas") // Retorna erro se o usuário não for encontrado
	}

	// Verifica se a senha fornecida corresponde ao hash armazenado
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("Senha incorreta para: %v", req.Email)
		return "", errors.New("credenciais inválidas") // Retorna erro se a senha for inválida
	}

	// Gera um token JWT para o usuário
	token, err := auth_utils.GenerateJWT(user.ID, user.Username, user.AccessLevel)
	if err != nil {
		logger.Error("Erro ao gerar token JWT: %v", err)
		return "", errors.New("erro ao gerar token") // Retorna erro se a geração do token falhar
	}

	return token, nil // Retorna o token gerado
}

// Register cria um novo usuário no banco de dados
//
// Parâmetros:
// - req: A requisição de registro contendo email, senha e nome do usuário.
//
// Retorno:
// - error: Retorna erro caso haja falha na criptografia da senha ou no registro do usuário.
//
// Detalhes:
// - A senha fornecida é criptografada e o usuário é registrado no banco de dados.
func Register(req RegisterRequest) error {
	// Criptografa a senha antes de armazená-la
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Erro ao criptografar senha: %v", err)
		return errors.New("erro ao registrar usuário") // Retorna erro se a criptografia falhar
	}

	// Cria o objeto de usuário com os dados fornecidos
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	// Registra o novo usuário no banco de dados
	if err := models.CreateNewUser(user); err != nil {
		logger.Error("Erro ao salvar usuário no banco: %v", err)
		return errors.New("erro ao registrar usuário") // Retorna erro caso falhe ao registrar o usuário
	}

	return nil // Retorna nil se o registro for bem-sucedido
}
