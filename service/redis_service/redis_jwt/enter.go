package redis_jwt

import (
	"blogx_server/global"
	"blogx_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type BlackListType int8

const (
	UserBlackListType   = 1 // 用户注销登录
	AdminBlackListType  = 2 // 管理员手动下线
	DeviceBlackListType = 3 // 其他设备把当前设备挤下来
)

func (b BlackListType) String() string {
	return fmt.Sprintf("%d", b)
}

func (b BlackListType) Msg() string {
	switch b {
	case UserBlackListType:
		return "已注销"
	case AdminBlackListType:
		return "禁止登录"
	case DeviceBlackListType:
		return "设备下线"
	}
	return "已注销"
}

func ParseBlackLitType(val string) BlackListType {
	switch val {
	case "1":
		return UserBlackListType
	case "2":
		return AdminBlackListType
	case "3":
		return DeviceBlackListType
	}
	return UserBlackListType
}

func TokenBlackList(token string, value BlackListType) {
	key := fmt.Sprintf("token_black_%s", token)

	claims, err := jwts.ParseToken(token)
	if err != nil || claims == nil {
		logrus.Errorf("token 解析失败 %s", err)
	}
	second := claims.ExpiresAt - time.Now().Unix()

	_, err = global.Redis.Set(key, value.String(), time.Duration(second)*time.Second).Result()
	if err != nil {
		logrus.Errorf("redis 添加黑名单失败", err)
	}
}

func HasTokenBlackList(token string) (blk BlackListType, ok bool) {
	key := fmt.Sprintf("token_black_%s", token)
	value, err := global.Redis.Get(key).Result()
	if err != nil {
		return
	}
	blk = ParseBlackLitType(value)
	return blk, true
}

func HasTokenBlackListByGin(c *gin.Context) (blk BlackListType, ok bool) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return HasTokenBlackList(token)
}
