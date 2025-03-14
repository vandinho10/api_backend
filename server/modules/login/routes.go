// pwd: /app/server/modules/login/router.go

package login

import (
	"api/logger"
	"api/server/modules/login/controllers"
	"api/server/modules/login/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes adiciona as rotas do módulo de autenticação ao roteador principal.
//
// Parâmetros:
//   - router: ponteiro para a instância do Gin Engine onde as rotas serão registradas.
func RegisterRoutes(router *gin.Engine) {
	// Grupo de rotas de autenticação
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - Auth"})
		})
		authGroup.POST("/login", controllers.AuthenticateUser)
		authGroup.POST("/logout", controllers.LogoutUser)
		authGroup.GET("/is_logged", middleware.AuthMiddleware(), controllers.IsLoggedIn)
	}

	// Grupo de rotas de administração de usuários (Nível Admin)
	usersGroup := authGroup.Group("/users").Use(middleware.AuthMiddleware())
	{
		usersGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - Users"})
		})
		usersGroup.POST("/register", controllers.AddNewUser)
		// usersGroup.GET("/", controllers.ListUsers)
		// usersGroup.PUT("/change_password", controllers.ChangePassword)
	}

	// Grupo de rotas para o usuário autenticado
	userGroup := authGroup.Group("/user").Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - User"})
		})
		// userGroup.GET("/me", controllers.GetUserDetails)
		// userGroup.GET("/", controllers.ListUsers)
		// userGroup.GET("/:id", controllers.GetUserByID)
		// userGroup.PUT("/:id", controllers.UpdateUser)
		// userGroup.DELETE("/:id", controllers.DeleteUser)
	}

	healthCheck()
}

// healthCheck registra uma mensagem de depuração indicando que o módulo de login está funcionando corretamente.
func healthCheck() {
	logger.Debug("Módulo Login está saudável.")
}

// init inicializa o módulo de login registrando uma mensagem de depuração no log.
func init() {
	logger.Debug("Módulo Login carregado.")
}

// FunctionTest é um exemplo de função para criar um usuário diretamente sem modificar o model.
//
// Esta função cria um usuário fictício e registra mensagens no log indicando o sucesso ou erro da operação.
// func FunctionTest() {
// 	err := models.CreateUser(models.User{
// 		Name:        "Teste User Default",
// 		Username:    "userdefault",
// 		Email:       "userdefault@example.com",
// 		Password:    "senha123",
// 		AccessLevel: 1,
// 	})

// 	if err != nil {
// 		logger.Error("Erro ao criar usuário: %v", err)
// 	} else {
// 		logger.Info("Usuário criado com sucesso!")
// 	}
// 	logger.Debug("Criando usuário.")
// }
