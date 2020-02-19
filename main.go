package main

import "ox-web/server"

// 一个类似gin的框架
// 路由使用前缀树
func main() {
	app := server.New()
	app.GET("/", func(context *server.Context) {
		context.JSON(200, server.M{
			"message": "go",
		})
	})
	app.Run(":8080")
}
