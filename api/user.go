package api

import (
	"fin_im/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserRegister(ctx *gin.Context) {
	logrus.Info("UserRegister")
	var userRegisterService service.UserRegisterService

	if err := ctx.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		ctx.JSON(200, res)
	} else {
		ctx.JSON(400, ErrorResponse(err))
		logrus.Info(err)
	}

}
