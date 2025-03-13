// pwd: /app/server/modules/login/router.go
package login

import (
	"api/logger"
	"api/server/modules/login/controllers"

	"api/server/modules/login/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes adiciona as rotas do módulo de autenticação
func RegisterRoutes(router *gin.Engine) {
	// Rotas de autenticação
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - Auth"})
		})
		authGroup.POST("/login", controllers.AuthenticateUser)
		authGroup.POST("/logout", controllers.LogoutUser)
		authGroup.GET("/is_logged", middleware.AuthMiddleware(), controllers.IsLoggedIn)
	}

	// Rotas de usuários(Nivel Admin)
	usersGroup := authGroup.Group("/users").Use(middleware.AuthMiddleware())
	{
		usersGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - Users"})
		})
		usersGroup.POST("/register", controllers.AddNewUser)
		// 	usersGroup.GET("/", controllers.ListUsers)
		// 	usersGroup.PUT("/change_password", controllers.ChangePassword)
	}

	// Rotas de usuário
	userGroup := authGroup.Group("/user").Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong - User"})
		})
		// usersGroup.GET("/me", middleware.AuthMiddleware(), controllers.GetUserDetails)
		// usersGroup.GET("/", controllers.ListUsers)
		// usersGroup.GET("/:id", controllers.GetUserByID)
		// usersGroup.PUT("/:id", controllers.UpdateUser)
		// usersGroup.DELETE("/:id", controllers.DeleteUser)
	}

	healthCheck()

}
func healthCheck() {
	logger.Debug("Módulo Login está saudável.")
}

func init() {
	logger.Debug("Módulo Login carregado.")
}

// func FunctionTest() {
// 	// Criar usuário diretamente na execução sem alterar o model
// 	err := models.CreateUser(models.User{
// 		Name:        "Teste User Default",
// 		Username:    "userdefault",
// 		Email:       "userdefault@example.com	",
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
