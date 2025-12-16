package log_service

import (
	"github.com/gin-gonic/gin"
	"minibar_server/core"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	claims, err := jwts.ParseTokenByGin(c)
	userID := uint(0)
	userName := ""
	if err == nil && claims != nil {
		userID = claims.UserID
		userName = claims.Username
	}

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录",
		Content:     "",
		UserID:      userID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		Username:    userName,
		Pwd:         "",
		LoginType:   loginType,
	})
}

func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, pwd string) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		Title:       "用户登录失败",
		Content:     msg,
		IP:          ip,
		Addr:        addr,
		LoginStatus: false,
		Username:    username,
		Pwd:         pwd,
		LoginType:   loginType,
	})
}
