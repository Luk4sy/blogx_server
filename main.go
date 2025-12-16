package main

import (
	"minibar_server/core"
	"minibar_server/flags"
	"minibar_server/global"
	"minibar_server/router"
)

func main() {
	flags.Parse()                   // 1. 解析命令行参数：flags
	global.Config = core.ReadConf() // 2. 加载配置文件：conf + core/init_conf.go
	core.InitLogrus()               // 3. 初始化日志：core/init_logrus.go + logs
	global.DB = core.InitDB()       // 4. 初始化数据库：core/init_db.go
	global.Redis = core.InitRedis() // 初始化redis
	global.ESClient = core.EsConnect()

	flags.Run()

	core.InitMysqlES()

	// 启动 web 程序
	router.Run() // 启动 HTTP 服务：r.Run(...)
}
