package main

import logger "ox-web/log"

// 一个类似gin的框架
// 路由使用前缀树
func main() {
	//app := core.New()
	//app.GET("/", func(context *core.Context) {
	//	context.JSON(200, core.M{
	//		"message": "go",
	//	})
	//})
	//app.Run(":8080")

	_ = logger.INFO().Str("greate", "spider man").Send()
}
