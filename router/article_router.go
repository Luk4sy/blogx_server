package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/api/article_api"
	"minibar_server/middleware"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.GET("article", middleware.BindQueryMiddleware[article_api.ArticleListRequest], app.ArticleListView)

}
