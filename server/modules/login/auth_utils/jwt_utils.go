// pwd: /app/server/modules/login/utils/jwt_utils.go
package auth_utils

import (
	"api/logger"
	"api/utils"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	AccessLevel int    `json:"access_level"`
	jwt.StandardClaims
}

var SecretKey = []byte(utils.GetEnv("JWT_SECRET"))

// GenerateJWT gera um novo token JWT para o ID de usuário fornecido.
func GenerateJWT(userID int, username string, accessLevel int) (string, error) {
	logger.Debug("Gerando um novo token JWT para o usuário ID=%v", userID)

	expirationTime, err := GetExpirationTime()
	if err != nil {
		logger.Error("Erro ao obter o tempo de expiração: %v", err)
		return "", errors.New("erro interno")
	}

	expirationDate := time.Now().Add(expirationTime)
	logger.Debug("Data de expiração do token: %v", expirationDate)
	claims := Claims{
		ID:          userID,
		Username:    username,
		AccessLevel: accessLevel,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		logger.Error("Erro ao assinar o token JWT: %v", err)
		return "", errors.New("erro ao gerar token JWT")
	}

	logger.Debug("Token gerado com sucesso para o usuário ID=%v", userID)
	return signedToken, nil
}

// HashAndCheckPassword gera um hash bcrypt da senha e compara com o hash fornecido.
func HashAndCheckPassword(password, hash string) (string, bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return string(hashedPassword), err == nil, nil
}

// GetExpirationTime recupera o tempo de expiração do JWT a partir da variável de ambiente "JWT_EXPIRE".
func GetExpirationTime() (time.Duration, error) {
	expirationStr := utils.GetEnv("JWT_EXPIRE")
	if expirationStr == "" {
		logger.Warn("Tempo de expiração não encontrado. Usando o valor padrão de 120 minutos")
		return 120 * time.Minute, nil
	}

	expirationTime, err := time.ParseDuration(expirationStr)
	if err != nil {
		return 0, fmt.Errorf("erro ao parsear o tempo de expiração: %v", err)
	}

	return expirationTime, nil
}

// ValidateJWT valida um token JWT fornecido.
func ValidateJWT(tokenString string) (bool, error) {
	logger.Debug("Validando o token JWT")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn("Método de assinatura inválido: %v", token.Header["alg"])
			return nil, errors.New("método de assinatura inválido")
		}
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		logger.Warn("Token inválido ou expirado: %v", err)
		return false, errors.New("token inválido")
	}

	logger.Debug("Token JWT é válido")
	return true, nil
}

// CalculateTokenExpirationTime calcula o tempo restante até a expiração de um token JWT.
func CalculateTokenExpirationTime(tokenString string) (time.Duration, error) {
	tokenString = RemoveBearerPrefix(tokenString)

	logger.Debug("Calculando o tempo restante até a expiração do token...")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn("Método de assinatura inválido: %v", token.Header["alg"])
			return nil, errors.New("método de assinatura inválido")
		}
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		logger.Warn("Token inválido ou expirado: %v", err)
		return 0, errors.New("token inválido ou expirado")
	}

	expirationTime := claims.ExpiresAt
	if expirationTime == 0 {
		logger.Warn("Tempo de expiração não encontrado no token")
		return 0, errors.New("tempo de expiração não encontrado")
	}

	expirationTimeUnix := time.Unix(expirationTime, 0)
	timeRemaining := time.Until(expirationTimeUnix)

	logger.Debug("Tempo restante até a expiração do token: %v", timeRemaining)

	if timeRemaining < 0 {
		logger.Warn("Token expirado")
		return 0, errors.New("token expirado")
	}

	return timeRemaining, nil
}

// RemoveBearerPrefix remove o prefixo "Bearer " de um cabeçalho de autenticação.
func RemoveBearerPrefix(authHeader string) string {
	logger.Debug("Conteúdo do cabeçalho de autenticação: %s", authHeader)
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return authHeader
}

// ValidateAndExtractClaims valida um token JWT e retorna os claims extraídos.
func ValidateAndExtractClaims(tokenString string) (*Claims, error) {
	logger.Debug("Token recebido: %v", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn("Método de assinatura inválido: %v", token.Header["alg"])
			return nil, errors.New("método de assinatura inválido")
		}
		return SecretKey, nil
	})

	logger.Debug("Validando o token JWT...")
	if err != nil {
		logger.Warn("Erro ao validar o token: %v", err)
		return nil, errors.New("token inválido ou expirado")
	}

	// Recupera os claims
	claims, ok := token.Claims.(jwt.MapClaims)
	logger.Debug("Token.Claims: %v", token.Claims.(jwt.MapClaims))

	logger.Debug("Claims extraídos do token: %v", claims)
	if !ok {
		logger.Warn("Falha ao extrair claims do token")
		return nil, errors.New("erro ao extrair claims do token")
	}

	// Converte jwt.MapClaims para Claims
	convertedClaims := &Claims{
		ID:          int(claims["id"].(float64)),
		Username:    claims["username"].(string),
		AccessLevel: int(claims["access_level"].(float64)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(claims["exp"].(float64)),
		},
	}
	return convertedClaims, nil
}
