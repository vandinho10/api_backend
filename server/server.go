// pwd: /app/server/server.go
package server

import (
	"api/logger"
	"api/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StartServer inicia o servidor HTTP e retorna a instância do servidor
func StartServer() (*http.Server, error) {
	logger.Info("Iniciando o servidor...")

	// Define a porta padrão para o servidor HTTP
	port_http := utils.GetEnv("PORT_HTTP")
	if port_http == "" {
		port_http = "80"
		logger.Warn("Variável de ambiente PORT_HTTP não definida. Usando valor padrão: %s", port_http)
	}

	// Define o endereço do servidor HTTP
	addr_http := fmt.Sprintf(":%s", port_http)
	logger.Info("Endereço do servidor HTTP definido para %s", addr_http)

	// Define o modo do Gin
	ginMode := utils.GetEnv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release" // Define um valor padrão
		logger.Warn("Variável de ambiente GIN_MODE não definida. Usando valor padrão: %s", ginMode)
	}
	gin.SetMode(ginMode)
	logger.Debug("Modo do Gin definido para %s", ginMode)

	// Função para permitir apenas o dominio e subdomínios específicados
	domainName := utils.GetEnv("DOMAIN_NAME")
	allowOrigins := func(origin string) bool {
		return strings.HasSuffix(origin, "."+domainName) || origin == "https://"+domainName
	}

	// Cria uma instância do Gin
	r := gin.New()

	// Definindo Use(configurações especificas) para o Gin
	r.Use(
		logger.LoggerMiddleware(),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc:  allowOrigins,
		}))

	// Define as rotas padrão: /health, /ping e /favicon.ico
	r.StaticFile("/favicon.ico", "./server/static/favicon.ico")
	r.GET("/health", HealthCheck)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Carrega os pacotes de rotas
	autoloadPackages(r)

	// Cria o servidor HTTP
	httpServer := &http.Server{
		Addr:    addr_http,
		Handler: r,
	}

	// Inicia o servidor HTTP em uma goroutine
	go func() {
		logger.Debug("Iniciando servidor HTTP na porta %s", addr_http)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erro ao iniciar o servidor HTTP: %v", err)
		}
	}()

	return httpServer, nil
}

// ShutdownServer realiza o shutdown gracioso do servidor HTTP
func ShutdownServer(ctx context.Context, server *http.Server) error {
	// Inicia o processo de desligamento do servidor com um contexto
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Erro ao desligar o servidor: %v", err)
		return fmt.Errorf("erro ao desligar o servidor: %w", err)
	}

	logger.Info("Servidor desligado com sucesso")
	return nil
}
