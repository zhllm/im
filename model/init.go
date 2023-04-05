package model

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Database(connstring string) {
	db, err := gorm.Open(mysql.Open(connstring), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	db.Logger.LogMode(logger.Info)
	if err != nil {
		fmt.Println("connect err", err.Error())
		panic(err)
	}
	DB = db
	logrus.Info("Mysql connect Successful!")

	if gin.Mode() == "release" {
		logrus.Info("release")
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	migration()
}
