package open

import (
	"go_gin/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	controllers.BaseController
}

func (con *UserController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "首页",
	})
}

func (con *UserController) Login(c *gin.Context) {
	con.Fail(c)
}
