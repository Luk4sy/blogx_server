package core

import (
	"github.com/sirupsen/logrus"
	"minibar_server/global"
	river "minibar_server/service/river_service"
)

func InitMysqlES() {
	if !global.Config.River.Enable {
		logrus.Infof("关闭 mysql 同步操作")
	}
	r, err := river.NewRiver()
	if err != nil {
		logrus.Fatal(err)
	}
	go r.Run()
}
