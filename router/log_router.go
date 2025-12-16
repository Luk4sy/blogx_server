package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
	"minibar_server/middleware"
)

func LogRouter(rr *gin.RouterGroup) {
	app := api.App.LogApi
	r := rr.Group("").Use(middleware.AdminMiddleware)
	r.GET("logs", app.LogListView)
	r.GET("logs/:id", app.LogReadView)
	r.DELETE("logs", app.LogRemoveView)

}
