// pwd: /app/server/modules/login/models/user_model.go
package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"api/db"
	"api/logger"
	"api/server/modules/login/auth_utils"

	"golang.org/x/crypto/bcrypt"
)

// User representa a estrutura do usuário no banco de dados
type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	AccessLevel int       `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
}

// LoginRequest representa os dados recebidos para login
type LoginRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse estrutura para a resposta de autenticação
type AuthResponse struct {
	Token         string `json:"token"`
	User          User   `json:"user"`
	TimeRemaining string `json:"time_remaining"`
}

// Authenticate verifica as credenciais do usuário e retorna um token JWT com as informações do usuário
func Authenticate(loginData LoginRequest) (*AuthResponse, error) {
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return nil, errors.New("erro interno de conexão")
	}
	defer dbConn.Close()

	var user User
	var query string
	var args []interface{}

	// Verifica se foi passado o Name, Username ou o Email
	if loginData.Username != "" {
		// Se passar o username, usar o username na consulta
		query = "SELECT id, username, email, password, access_level FROM users WHERE username = ? LIMIT 1"
		args = append(args, loginData.Username)
	} else {
		return nil, errors.New("username é necessário para o login")
	}

	// Executa a consulta
	err = dbConn.QueryRow(query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AccessLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Usuário não encontrado: %v", loginData.Name)
			return nil, errors.New("usuário ou senha inválidos")
		}
		logger.Error("Erro ao buscar usuário: %v", err)
		return nil, errors.New("erro interno")
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		logger.Warn("Tentativa de login com senha incorreta para usuário %v", loginData.Username)
		return nil, errors.New("usuário ou senha inválidos")
	}

	// Gera o token JWT
	token, err := auth_utils.GenerateJWT(user.ID, user.Username, user.AccessLevel)
	if err != nil {
		logger.Error("Erro ao gerar token JWT: %v", err)
		return nil, errors.New("erro interno ao gerar token")
	}

	// Calcular o tempo restante até o vencimento do token (exemplo: 3600 segundos)
	// Pode ser configurado de acordo com a lógica de expiração do JWT
	token = "Bearer " + token
	timeRemaining, err := auth_utils.CalculateTokenExpirationTime(token) // Exemplo de tempo restante em segundos (1 hora)
	if err != nil {
		logger.Error("Erro ao calcular o tempo de expiração do token: %v", err)
		return nil, errors.New("erro interno ao calcular o tempo de expiração do token")
	}

	// Retorna o token e o usuário com o tempo restante
	authResponse := &AuthResponse{
		Token:         token,
		User:          user,
		TimeRemaining: timeRemaining.String(),
	}

	return authResponse, nil
}

// GetUserByEmail busca um usuário pelo e-mail no banco de dados
func GetUserByEmail(email string) (*User, error) {
	// Estabelece a conexão com o banco de dados
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return nil, errors.New("erro interno de conexão")
	}
	defer dbConn.Close()

	var user User
	query := "SELECT id, name, email, password FROM users WHERE email = ? LIMIT 1"
	err = dbConn.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Usuário não encontrado: %v", email)
			return nil, errors.New("usuário não encontrado")
		}
		logger.Error("Erro ao buscar usuário por email: %v", err)
		return nil, errors.New("erro interno")
	}

	return &user, nil
}

// CreateNewUser insere um novo usuário no banco de dados
func CreateNewUser(user User) error {
	logger.Debug("Criando novo usuário: %v", user.Username)
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return errors.New("erro interno de conexão")
	}
	defer dbConn.Close()

	// Hash da senha antes de armazenar
	logger.Debug("Gerando hash da senha para o usuário %v", user.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	logger.Debug("Hash gerado: %v | Base: %v", string(hashedPassword), user.Password)

	if err != nil {
		logger.Error("Erro ao gerar hash da senha: %v", err)
		return errors.New("erro ao processar senha")
	}

	query := "INSERT INTO users (name, username, email, password, access_level) VALUES (?, ?, ?, ?, ?)"
	_, err = dbConn.Exec(query, user.Name, user.Username, user.Email, string(hashedPassword), user.AccessLevel) // Converte o hash para string
	if err != nil {
		logger.Error("Erro ao criar usuário: %v", err)
		return errors.New("erro ao registrar usuário")
	}

	return nil
}

// CheckUserExists verifica se o nome de usuário já existe no banco
func CheckUserExists(username string) error {
	// Estabelece a conexão com o banco de dados
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return errors.New("erro interno de conexão")
	}
	defer dbConn.Close()

	var user User
	query := "SELECT id FROM users WHERE username = ? LIMIT 1"
	err = dbConn.QueryRow(query, username).Scan(&user.ID)
	logger.Debug(("Resultado da consulta: %v"), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Usuário não encontrado: %v", username)
			return nil // Usuário não encontrado, pode criar um novo
		}
		logger.Error("Erro ao verificar usuário: %v", err)
		return errors.New("erro interno")
	}

	return fmt.Errorf("usuário já existe")
}

// AddTokenToBlacklist adiciona um token à blacklist de tokens JWT no banco de dados.
//
// Parâmetros:
// - tokenString: Token JWT a ser invalidado.
// - expiresAt: Data e hora de expiração do token.
//
// Retorno:
// - Retorna erro caso ocorra falha na inserção.
func AddTokenToBlacklist(tokenString string, expiresAt time.Time) error {
	// Estabelece conexão com o banco de dados
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return errors.New("erro interno de conexão")
	}
	defer dbConn.Close()

	query := "INSERT INTO jwt_blacklist (token, expires_at) VALUES (?, ?)"
	_, err = dbConn.Exec(query, tokenString, expiresAt)
	if err != nil {
		logger.Error("Erro ao adicionar token à blacklist: %v", err)
		return errors.New("erro interno ao adicionar token à blacklist")
	}

	return nil
}

// IsTokenBlacklisted verifica se um token está na blacklist.
// Retorna true se o token estiver na blacklist, false caso contrário.
// Em caso de erro, assume-se que o token é válido (retorna false).
func IsTokenBlacklisted(tokenString string) bool {
	dbConn, err := db.DbConnection()
	if err != nil {
		logger.Error("Erro ao conectar ao banco de dados: %v", err)
		return false // Assume token válido em caso de erro
	}
	defer dbConn.Close()

	var count int
	err = dbConn.QueryRow("SELECT COUNT(*) FROM jwt_blacklist WHERE token = ?", tokenString).Scan(&count)
	if err != nil {
		logger.Error("Erro ao verificar blacklist: %v", err)
		return false // Assume token válido em caso de erro
	}

	return count > 0
}
