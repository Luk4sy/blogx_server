package user_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/utils/pwd"
)

type ResetPwdRequest struct {
	Pwd string `json:"pwd" binding:"required"`
}

func (UserApi) ResetPwdView(c *gin.Context) {
	var cr ResetPwdRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	var user models.UserModel
	err = global.DB.Take(&user, "email = ?", email).Error
	if err != nil {
		res.FailWithMsg("不存在的用户", c)
		return
	}

	if user.RegisterSource != enum.RegisterEmailSourceType {
		res.FailWithMsg("非邮箱注册用户，无重置密码权限", c)
	}

	hashPwd, _ := pwd.GenerateFromPassword(cr.Pwd)
	global.DB.Model(&user).Update("password", hashPwd)
	res.OkWithMsg("重置密码成功", c)

}
