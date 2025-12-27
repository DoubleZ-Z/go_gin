package routers

import (
	"go_gin/controllers/open"
	"go_gin/interceptor"

	"github.com/gin-gonic/gin"
)

func OpenRouters(engine *gin.Engine) {
	openGroup := engine.Group("/api/open")
	{
		userController := open.UserController{}
		openGroup.GET("/", userController.Home)
		openGroup.GET("/login", interceptor.TokenAuth, userController.Login)
	}
}
