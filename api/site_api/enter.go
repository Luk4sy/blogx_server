package site_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

type SiteInfoRequest struct {
	Name string `uri:"name"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	var cr SiteInfoRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if cr.Name == "site" {
		res.OkWithData(global.Config.Site, c)
		return
	}

	// 判断角色是不是管理员
	middleware.AdminMiddleware(c)

	_, ok := c.Get("claims")
	if !ok {
		return
	}

	var data any
	switch cr.Name {
	case "email":
		rep := global.Config.Email
		rep.AuthCode = "******"
		data = rep
	case "qq":
		rep := global.Config.QQ
		rep.AppKey = "******"
		data = rep
	case "qiNiu":
		rep := global.Config.QiNiu
		rep.SecretKey = "******"
		data = rep
	case "ai":
		rep := global.Config.Ai
		rep.SecretKey = "******"
		data = rep
	default:
		res.FailWithMsg("不存在的配置", c)
		return
	}
	res.OkWithData(data, c)
	return
}

func (SiteApi) SiteInfoQQView(c *gin.Context) {
	res.OkWithData(global.Config.QQ.Url(), c)
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	var cr SiteUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithMsg("更新成功", c)
	return
}
