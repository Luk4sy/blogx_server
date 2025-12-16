package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/middleware"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banner", app.BannerListView)
	r.POST("banner", middleware.AdminMiddleware, app.BannerCreateView)
	r.DELETE("banner", middleware.AdminMiddleware, app.BannerRemoveView)
	r.PUT("banner/:id", middleware.AdminMiddleware, app.BannerUpdateView)

}
