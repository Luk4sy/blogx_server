package user_service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"minibar_server/core"
	"minibar_server/global"
	"minibar_server/models"
)

func (u UserService) UserLogin(c *gin.Context) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)
	ua := c.GetHeader("User-Agent")
	err := global.DB.Create(&models.UserLoginModel{
		UserID: u.userModel.ID,
		IP:     ip,
		Addr:   addr,
		UA:     ua,
	}).Error
	if err != nil {
		logrus.Errorf("用户登录日志写入失败 %s", err)
	}
}
