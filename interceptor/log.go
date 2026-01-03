package interceptor

import (
	"fmt"
	"go_gin/util"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	cp := c.Copy()
	go fmt.Println(util.UnixToTimeString(time.Now().Unix()), cp.Request.Host, cp.Request.RequestURI, "这里异步记录了一个日志")
}
