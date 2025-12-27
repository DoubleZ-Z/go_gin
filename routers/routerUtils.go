package routers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeCost(c *gin.Context) {
	now := time.Now()
	c.Next()
	cost := time.Since(now)
	fmt.Println(c.Request.Host, c.Request.RequestURI, "时间消耗:", cost)
}
