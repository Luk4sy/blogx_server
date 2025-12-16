package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/middleware"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.POST("images", middleware.AuthMiddleware, app.ImageUploadView)
	r.POST("images/qiniu", middleware.AuthMiddleware, app.QiNiuGenToken)
	r.GET("images", middleware.AdminMiddleware, app.ImageListView)
	r.DELETE("images", middleware.AdminMiddleware, app.ImageRemoveView)

}
