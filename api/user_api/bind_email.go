package user_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/utils/jwts"
)

func (UserApi) BindEmailView(c *gin.Context) {

	//TODO: 完善邮箱绑定
	// 关于安全性，邮箱绑定要考虑的很多
	// 1. 绑定的邮箱是否有人使用？
	// 2. 新邮箱替换旧邮箱？ 旧邮箱是否能用？
	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("不存在的用户", c)
		return
	}
	global.DB.Model(&user).Update("email", email)
	res.OkWithMsg("邮箱绑定成功", c)
}
