package interceptor

import (
	"github.com/gin-gonic/gin"
)

func TokenAuth(c *gin.Context) {
	/*token := c.Request.Header.Get("Token")
	if token == "" {
		c.Abort()
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "请登录",
	})*/
}
