package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/middleware"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("site/qq_url", app.SiteInfoQQView)
	r.GET("site/:name", app.SiteInfoView)
	r.PUT("site/:name", middleware.AdminMiddleware, app.SiteUpdateView)

}
