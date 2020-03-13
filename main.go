package main

import (
	"log"
	"ox-web/web"
	"ox-web/websocket"
)

// 一个类似gin的框架
// 路由使用前缀树
func main() {
	app := ox.New()
	app.GET("/", func(ctx *ox.Context) {
		conn, err := websocket.Upgrade(ctx.Writer, ctx.Req)
		if err != nil {
			log.Println(err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte("你好")); err != nil {
			log.Println(err)
		}
		for {
			fin, op, body, err := conn.Read()
			log.Println(fin, op, body, err)
		}
	})
	_ = app.Run(":8080")
}

