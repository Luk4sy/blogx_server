package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/redis_service/redis_jwt"
	"fmt"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Redis = core.InitRedis()

	//token, err := jwts.GetToken(jwts.Claims{
	//	UserID: 2,
	//	Role:   1,
	//})
	//fmt.Println(token, err)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjIsInVzZXJuYW1lIjoiIiwicm9sZSI6MSwiZXhwIjoxNzYzNDc4MjI2LCJpc3MiOiJMdWthc3kifQ.a0U4A1Sf7zQYwWBV0RUV0GhPZDpq-483m-6jMSRirXY"
	redis_jwt.TokenBlackList(token, redis_jwt.UserBlackListType)
	blk, ok := redis_jwt.HasTokenBlackList(token)
	fmt.Println(blk, ok)
}
