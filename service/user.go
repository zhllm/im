package service

import (
	"fin_im/model"
	"fin_im/serializer"

	"github.com/sirupsen/logrus"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	var count int64 = 0
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).Count(&count).First(&user)
	logrus.Info("get register count ", count, service.UserName, "  ", user.UserName, " xxxxx")
	if count != 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "用户名已经存在",
		}
	}

	user = model.User{
		UserName: service.UserName,
	}
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "加密错误",
		}
	}

	model.DB.Create(&user)
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功",
	}
}
