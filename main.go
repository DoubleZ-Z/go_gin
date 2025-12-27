package main

import (
	"go_gin/interceptor"
	"go_gin/routers"

	"github.com/gin-gonic/gin"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	engine := gin.Default()
	engine.Use(interceptor.Log, routers.TimeCost)
	routers.OpenRouters(engine)
	routers.TokenRouters(engine)
	err := engine.Run()
	if err != nil {
		panic(err)
		return
	}
}
