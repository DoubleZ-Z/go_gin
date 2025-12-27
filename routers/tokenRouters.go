package routers

import (
	"go_gin/controllers/token"
	"go_gin/interceptor"

	"github.com/gin-gonic/gin"
)

func TokenRouters(engine *gin.Engine) {
	group := engine.Group("/api/token", interceptor.TokenAuth)
	{
		controller := token.FileController{}
		group.POST("/file/upload", controller.UploadFile)
	}
}
