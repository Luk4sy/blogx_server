package router

import (
	"github.com/gin-gonic/gin"
	"minibar_server/api"
)

func CaptchaRouter(r *gin.RouterGroup) {
	app := api.App.CaptchaApi
	r.GET("captcha", app.CaptchaView)

}
