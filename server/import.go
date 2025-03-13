// pwd: /app/server/import.go
package server

import (
	"api/logger"
	"api/server/modules/finance"
	"api/server/modules/login"
	"api/server/modules/ppr"

	"github.com/gin-gonic/gin"
)

// Mapeia os módulos e suas respectivas funções de registro de rotas
// Isso permite adicionar novos módulos dinamicamente, sem a necessidade de alterações no código principal.
// Inlusãp manual somente aqui e no import Acima
var moduleRegistry = map[string]func(*gin.Engine){
	"finance": finance.RegisterRoutes, // Módulo de finanças
	"login":   login.RegisterRoutes,   // Módulo de login
	"ppr":     ppr.RegisterRoutes,     // Módulo de PPR (Plano de Participação nos Resultados)
}

// autoloadPackages registra automaticamente as rotas de cada módulo
//
// Essa função percorre o registro de módulos e chama a função de registro de rotas
// de cada módulo, permitindo a adição dinâmica de novos módulos sem alterações
// diretas no código principal.
//
// Parâmetros:
//   - router (*gin.Engine): Instância do roteador Gin onde as rotas serão registradas.
func autoloadPackages(router *gin.Engine) {
	// Itera sobre o mapa de módulos e registra as rotas de cada módulo
	for name, registerFunc := range moduleRegistry {
		// Loga o nome do módulo que está sendo registrado
		logger.Debug("Registrando rotas para o módulo: %s", name)
		// Chama a função de registro de rotas do módulo
		registerFunc(router)
	}
}
