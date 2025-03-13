// pwd: /app/server/import.go
package server

import (
	"api/logger"
	"api/server/modules/finance"
	"api/server/modules/login"
	"api/server/modules/ppr"

	"github.com/gin-gonic/gin"
)

// Mapeia os módulos para facilitar a adição dinâmica
var moduleRegistry = map[string]func(*gin.Engine){
	"finance": finance.RegisterRoutes,
	"login":   login.RegisterRoutes,
	"ppr":     ppr.RegisterRoutes,
}

// Registra automaticamente as rotas de cada módulo
func autoloadPackages(router *gin.Engine) {
	for name, registerFunc := range moduleRegistry {
		logger.Debug("Registrando rotas para o módulo: %s", name)
		registerFunc(router)
	}
}
