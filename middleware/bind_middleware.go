package middleware

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
)

func BindJsonMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	return
}

func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	return
}

func GetBind[T any](c *gin.Context) (cr T) {
	return c.MustGet("request").(T)
}
