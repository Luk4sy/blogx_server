package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/email_service"
	"blogx_server/utils/email_store"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

type SendEmailRequest struct {
	Type  int8   `json:"type" binding:"oneof=1 2 3"` // 1 注册 2 重置密码 3 绑定邮箱
	Email string `json:"email" binding:"required"`
}

type SendEmailResponse struct {
	EmailID string `json:"emailId"`
}

func (UserApi) SendEmailView(c *gin.Context) {
	var cr SendEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	code := base64Captcha.RandText(4, "012356789")
	id := base64Captcha.RandomId()

	switch cr.Type {
	case 1:
		// 查邮箱是否不存在
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已存在", c)
			return
		}
		err = email_service.SendRegisterCode(cr.Email, code)
	case 2:
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err != nil {
			res.FailWithMsg("该邮箱不存在", c)
			return
		}
		// 判断是否为邮箱注册
		if user.RegisterSource != enum.RegisterEmailSourceType {
			res.FailWithMsg("非邮箱注册用户，无重置密码权限", c)
		}
		err = email_service.SendResetPwdCode(cr.Email, code)
	case 3:
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已存在", c)
			return
		}
		err = email_service.SendBindEmailCode(cr.Email, code)
	}
	if err != nil {
		logrus.Errorf("邮件发送失败 %s", err)
		res.FailWithMsg("邮件发送失败", c)
	}
	// 使用 email_store.Set 存储验证码及时间戳
	email_store.Set(id, cr.Email, code)
	res.OkWithData(SendEmailResponse{
		EmailID: id,
	}, c)
}
