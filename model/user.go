package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
}

const (
	PassWordCast = 12 // 密码加密难度
)

func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCast)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
