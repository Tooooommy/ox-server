package main

import (
	"ox-web/core"
)

// 一个类似gin的框架
// 路由使用前缀树
func main() {
	app := core.New()
	app.GET("/", func(context *core.Context) {
		context.JSON(200, core.M{
			"message": "go",
		})
	})
	_ = app.Run(":8080")
}

//func Upgrade(writer http.ResponseWriter, request *http.Request) error {
//	if request.Method != http.MethodGet {
//		return errors.New("bad request, method not allowed")
//	}
//	if request.Header.Get("Sec-Websocket-Version") != 13 {
//
//	}
//}
