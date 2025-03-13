// pwd: /app/db/db.go

package db

import (
	"database/sql"
	"fmt"

	"api/logger"
	"api/utils"

	_ "github.com/go-sql-driver/mysql" // Importa o driver para registrar o driver MySQL.
)

// DbConnection estabelece uma conexão com o banco de dados MySQL/MariaDB.
//
// Esta função recupera os parâmetros de conexão das variáveis de ambiente,
// valida essas informações e tenta abrir uma conexão com o banco de dados.
//
// Retorna:
//   - *sql.DB: Ponteiro para a conexão com o banco de dados, se bem-sucedida.
//   - error: Erro detalhado em caso de falha ao conectar.
func DbConnection() (*sql.DB, error) {
	// Obtém os parâmetros de conexão do banco de dados a partir das variáveis de ambiente.
	dbUser, dbPass, dbHost, dbPort, dbName := getDBConfig()

	// Valida se todos os parâmetros obrigatórios foram fornecidos.
	if err := validateDBConfig(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
		return nil, err
	}

	// Constrói a string DSN (Data Source Name) para a conexão com o banco de dados.
	dsn := buildDSN(dbUser, dbPass, dbHost, dbPort, dbName)

	// Tenta abrir uma conexão com o banco de dados utilizando a DSN.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com o banco de dados: %w", err) // Adiciona contexto ao erro.
	}

	// Verifica se a conexão com o banco de dados foi estabelecida com sucesso.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err) // Adiciona contexto ao erro.
	}

	return db, nil // Retorna a conexão ativa com o banco de dados.
}

// getDBConfig recupera os parâmetros da conexão com o banco de dados a partir das variáveis de ambiente.
//
// Retorna:
//   - string: Usuário do banco de dados.
//   - string: Senha do banco de dados.
//   - string: Host do banco de dados.
//   - string: Porta de conexão do banco de dados.
//   - string: Nome do banco de dados.
func getDBConfig() (string, string, string, string, string) {
	return utils.GetEnv("DB_USER"), utils.GetEnv("DB_PASS"), utils.GetEnv("DB_HOST"), utils.GetEnv("DB_PORT"), utils.GetEnv("DB_NAME")
}

// buildDSN constrói a string DSN (Data Source Name) para a conexão com o banco de dados.
//
// A função também registra um log contendo a DSN de forma segura (omitindo a senha).
//
// Parâmetros:
//   - user (string): Usuário do banco de dados.
//   - password (string): Senha do banco de dados.
//   - host (string): Host do banco de dados.
//   - port (string): Porta de conexão do banco de dados.
//   - dbName (string): Nome do banco de dados.
//
// Retorna:
//   - string: DSN formatada para a conexão com o banco de dados.
func buildDSN(user, password, host, port, dbName string) string {
	// Monta a DSN sem a senha para log seguro.
	safeDSN := fmt.Sprintf("%s:***@tcp(%s:%s)/%s", user, host, port, dbName)
	logger.Debug("DSN: %s", safeDSN)

	// Retorna a DSN completa para a conexão com o banco de dados.
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

// validateDBConfig verifica se todas as variáveis essenciais para a conexão com o banco de dados foram fornecidas.
//
// Parâmetros:
//   - params (string...): Lista de parâmetros obrigatórios para validação.
//
// Retorna:
//   - error: Erro indicando a ausência de um ou mais parâmetros obrigatórios.
func validateDBConfig(params ...string) error {
	for _, param := range params {
		if param == "" {
			return fmt.Errorf("variável de ambiente obrigatória do banco de dados ausente")
		}
	}
	return nil
}
