package router

import (
	"blogx_server/global"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default() // 1️⃣ 创建一个默认的 Gin 路由引擎（带 Logger 和 Recovery 中间件）

	nr := r.Group("/api") // 2️⃣ 创建一个路由分组，所有接口都将以 "/api" 开头
	SiteRouter(nr)        // 3️⃣ 注册具体的路由模块，比如网站信息接口

	addr := global.Config.System.Addr() // 4️⃣ 从全局配置读取服务启动地址（如 :8080）
	r.Run(addr)                         // 5️⃣ 启动 HTTP 服务
}
