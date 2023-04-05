package router

import (
	"fin_im/api"
	"fin_im/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())

	v1 := r.Group("/")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, "Success")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}
	return r
}
