package router

import (
	"blogx_server/global"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode) // 配置路由模式
	r := gin.Default()                        // 1️⃣ 创建一个默认的 Gin 路由引擎（带 Logger 和 Recovery 中间件）

	r.Static("/uploads", "uploads") // 静态路由

	nr := r.Group("/api") // 2️⃣ 创建一个路由分组，所有接口都将以 "/api" 开头

	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr) // 注册具体的路由模块，站点模块
	LogRouter(nr)
	ImageRouter(nr)
	BannerRouter(nr)
	CaptchaRouter(nr)
	UserRouter(nr)

	addr := global.Config.System.Addr() // 4️⃣ 从全局配置读取服务启动地址（如 :8080）
	r.Run(addr)                         // 5️⃣ 启动 HTTP 服务
}
