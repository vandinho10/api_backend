// pwd: /app/utils/utils.go
package utils

import (
	"api/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var once sync.Once

// loadEnv garante que as variáveis de ambiente sejam carregadas apenas uma vez.
func loadEnv() {
	once.Do(func() {
		err := godotenv.Load("/app/.env")
		if err != nil {
			logger.Warn("Aviso: Arquivo .env não encontrado ou não pôde ser carregado")
		} else {
			logger.Info("Variáveis de ambiente carregadas com sucesso")
		}
	})
}

// GetEnv retorna uma variável de ambiente, garantindo que as variáveis sejam carregadas antes.
func GetEnv(key string) string {
	loadEnv()
	return os.Getenv(key)
}

// GracefulShutdown captura sinais de término e realiza o shutdown gracioso
func GracefulShutdown(server *http.Server) {
	// Cria um canal para capturar sinais de término (como SIGINT ou SIGTERM)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Aguarda até receber um sinal de término
	sig := <-sigs
	logger.Info("Recebido sinal de término: %s", sig)

	// Aqui você pode iniciar o processo de shutdown gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Erro ao encerrar o servidor: %v", err)
	}

	logger.Info("Servidor encerrado com sucesso")
}
