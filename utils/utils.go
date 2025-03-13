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

// loadEnv carrega as variáveis de ambiente do arquivo .env.
//
// Garante que as variáveis sejam carregadas apenas uma vez, utilizando sync.Once.
//
// Não retorna valores.
func loadEnv() {
	once.Do(func() {
		// Tenta carregar o arquivo .env
		err := godotenv.Load("/app/.env")
		if err != nil {
			// Se o arquivo .env não for encontrado ou não puder ser carregado, emite um aviso.
			logger.Warn("Aviso: Arquivo .env não encontrado ou não pôde ser carregado")
		} else {
			// Caso as variáveis de ambiente sejam carregadas com sucesso, emite uma informação.
			logger.Info("Variáveis de ambiente carregadas com sucesso")
		}
	})
}

// GetEnv retorna o valor de uma variável de ambiente.
//
// Certifica-se de que as variáveis de ambiente foram carregadas antes de retornar o valor.
//
// Parâmetros:
//   - key (string): Nome da variável de ambiente.
//
// Retorna:
//   - string: Valor da variável de ambiente ou uma string vazia caso a variável não exista.
func GetEnv(key string) string {
	loadEnv()
	return os.Getenv(key)
}

// GracefulShutdown realiza o shutdown gracioso do servidor HTTP.
//
// Espera até receber um sinal de término (SIGINT, SIGTERM, ou SIGQUIT) e encerra o servidor com um timeout.
//
// Parâmetros:
//   - server (*http.Server): O servidor HTTP que será encerrado.
//
// Não retorna valores.
func GracefulShutdown(server *http.Server) {
	// Cria um canal para capturar sinais de término (como SIGINT ou SIGTERM)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Aguarda até receber um sinal de término
	sig := <-sigs
	// Emite um log informando que o sinal de término foi recebido.
	logger.Info("Recebido sinal de término: %s", sig)

	// Inicia o processo de shutdown gracioso com um contexto de timeout de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tenta realizar o shutdown do servidor dentro do tempo limite
	if err := server.Shutdown(ctx); err != nil {
		// Se houver erro durante o processo de shutdown, emite um log de erro fatal
		logger.Fatal("Erro ao encerrar o servidor: %v", err)
	}

	// Emite um log informando que o servidor foi encerrado com sucesso
	logger.Info("Servidor encerrado com sucesso")
}
