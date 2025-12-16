package core

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"minibar_server/global"
)

func InitRedis() *redis.Client {
	r := global.Config.Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr:     r.Addr,     // 不写默认就是这个
		Password: r.Password, // 密码
		DB:       r.DB,       // 默认是0
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		logrus.Fatalf("redis 连接失败 %s", err)
	}
	logrus.Info("redis 连接成功 ")
	return redisDB
}
