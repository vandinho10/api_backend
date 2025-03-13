// pwd: /app/db/db.go

package db

import (
	"database/sql"
	"fmt"

	"api/logger"
	"api/utils"

	_ "github.com/go-sql-driver/mysql" // Importa o driver para registrar o driver MySQL
)

// DbConnection establishes a connection to the MySQL/MariaDB database.
// It retrieves the database connection parameters from environment variables.
// Returns a pointer to the database connection and an error if the connection fails.
func DbConnection() (*sql.DB, error) {
	// Fetch the database connection parameters from environment variables.
	dbUser, dbPass, dbHost, dbPort, dbName := getDBConfig()

	// Validate the database configuration.
	if err := validateDBConfig(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
		return nil, err
	}

	// Construct the Data Source Name (DSN) for connecting to the database.
	dsn := buildDSN(dbUser, dbPass, dbHost, dbPort, dbName)

	// Attempt to open a connection to the database using the DSN.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err) // Wrap the error for more context.
	}

	// Verify the connection to the database.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err) // Wrap the error for more context.
	}

	return db, nil // Return the database connection if successful.
}

// getDBConfig retrieves database configuration from environment variables.
func getDBConfig() (string, string, string, string, string) {
	return utils.GetEnv("DB_USER"), utils.GetEnv("DB_PASS"), utils.GetEnv("DB_HOST"), utils.GetEnv("DB_PORT"), utils.GetEnv("DB_NAME")
}

// buildDSN constructs the Data Source Name (DSN) for the database connection.
func buildDSN(user, password, host, port, dbName string) string {
	safeDSN := fmt.Sprintf("%s:***@tcp(%s:%s)/%s", user, host, port, dbName)
	logger.Debug("DSN: %s", safeDSN)

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}

func validateDBConfig(params ...string) error {
	for _, param := range params {
		if param == "" {
			return fmt.Errorf("missing required database environment variable")
		}
	}
	return nil
}
