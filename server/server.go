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

// StartServer inicia o servidor HTTP e retorna a instância do servidor.
//
// Configura e inicia o servidor HTTP, incluindo configuração de CORS e log.
// Caso a variável de ambiente PORT_HTTP não esteja definida, usa a porta padrão "80".
// O modo do Gin é definido a partir da variável de ambiente GIN_MODE, ou "release" por padrão.
//
// Retorna:
//   - *http.Server: Instância do servidor HTTP.
//   - error: Erro, caso ocorra durante a execução.
func StartServer() (*http.Server, error) {
	// Loga a inicialização do servidor
	logger.Info("Iniciando o servidor...")

	// Define a porta HTTP a partir da variável de ambiente, ou usa "80" por padrão
	port_http := utils.GetEnv("PORT_HTTP")
	if port_http == "" {
		port_http = "80"
		logger.Warn("Variável de ambiente PORT_HTTP não definida. Usando valor padrão: %s", port_http)
	}

	// Define o endereço completo do servidor HTTP
	addr_http := fmt.Sprintf(":%s", port_http)
	logger.Info("Endereço do servidor HTTP definido para %s", addr_http)

	// Define o modo de operação do Gin (ex: "release" ou "debug")
	ginMode := utils.GetEnv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release" // Define um valor padrão
		logger.Warn("Variável de ambiente GIN_MODE não definida. Usando valor padrão: %s", ginMode)
	}
	// Define o modo do Gin
	gin.SetMode(ginMode)
	logger.Debug("Modo do Gin definido para %s", ginMode)

	// Função para permitir apenas o domínio e subdomínios especificados
	domainName := utils.GetEnv("DOMAIN_NAME")
	allowOrigins := func(origin string) bool {
		// Permite apenas origens que correspondem ao domínio configurado
		return strings.HasSuffix(origin, "."+domainName) || origin == "https://"+domainName
	}

	// Cria uma nova instância do Gin
	r := gin.New()

	// Aplica middlewares: log e CORS
	r.Use(
		logger.LoggerMiddleware(),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc:  allowOrigins,
		}))

	// Define rotas padrão do servidor (health check, ping e favicon)
	r.StaticFile("/favicon.ico", "./server/static/favicon.ico")
	r.GET("/health", HealthCheck)
	r.GET("/ping", func(c *gin.Context) {
		// Retorna uma resposta "pong" para a rota /ping
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Carrega pacotes de rotas adicionais
	autoloadPackages(r)

	// Cria a instância do servidor HTTP
	httpServer := &http.Server{
		Addr:    addr_http,
		Handler: r,
	}

	// Inicia o servidor HTTP em uma goroutine para não bloquear o fluxo
	go func() {
		logger.Debug("Iniciando servidor HTTP na porta %s", addr_http)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Caso ocorra algum erro ao iniciar o servidor, loga como erro fatal
			logger.Fatal("Erro ao iniciar o servidor HTTP: %v", err)
		}
	}()

	// Retorna a instância do servidor HTTP
	return httpServer, nil
}

// ShutdownServer realiza o shutdown gracioso do servidor HTTP.
//
// Finaliza o servidor de maneira controlada, permitindo que as conexões ativas sejam completadas dentro de um tempo limite.
// Caso ocorra algum erro durante o desligamento, ele será retornado.
//
// Parâmetros:
//   - ctx (context.Context): Contexto com o timeout para o desligamento do servidor.
//   - server (*http.Server): Instância do servidor HTTP a ser desligado.
//
// Retorna:
//   - error: Erro, caso ocorra durante o processo de desligamento.
func ShutdownServer(ctx context.Context, server *http.Server) error {
	// Inicia o processo de desligamento gracioso com o contexto
	if err := server.Shutdown(ctx); err != nil {
		// Se ocorrer erro ao desligar o servidor, loga o erro
		logger.Error("Erro ao desligar o servidor: %v", err)
		return fmt.Errorf("erro ao desligar o servidor: %w", err)
	}

	// Loga a finalização do servidor com sucesso
	logger.Info("Servidor desligado com sucesso")
	return nil
}
