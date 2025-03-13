// pwd: /app/main.go
package main

import (
	"api/logger"
	"api/server"
	"api/utils"
)

func main() {
	// Inicializa o servidor e captura a instância do http.Server
	httpServer, err := server.StartServer()
	if err != nil {
		logger.Fatal("Erro ao iniciar o servidor: %v", err)
	}

	// Chama a função de utilitário para capturar sinais e fazer o shutdown gracioso
	utils.GracefulShutdown(httpServer)
}
