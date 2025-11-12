package site_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	// TODO：之后修改站点的信息
	fmt.Println("1")
	c.JSON(200, gin.H{"code": 0, "msg": "站点信息"})
	return
}
