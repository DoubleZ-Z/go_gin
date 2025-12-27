package interceptor

import (
	"fmt"
	"go_gin/models"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	cp := c.Copy()
	go fmt.Println(models.UnixToTime(time.Now().Unix()), cp.Request.Host, cp.Request.RequestURI, "这里异步记录了一个日志")
}
