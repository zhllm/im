package main

import (
	"fin_im/conf"
	"fin_im/router"
	"fin_im/service"

	"github.com/gin-gonic/gin"
)

func main() {
	println("Hello, WebAssembly!")
	conf.Init()
	gin.SetMode(gin.DebugMode)

	go service.Manager.Start()

	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
