package captcha_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

type CaptchaApi struct {
}

type CaptchaResponse struct {
	CaptchaID string `json:"captchaID"`
	Captcha   string `json:"captcha"`
}

func (CaptchaApi) CaptchaView(c *gin.Context) {
	// 使用内置的数字验证码驱动，先确认功能正常
	driver := base64Captcha.NewDriverDigit(
		60,  // 高度
		200, // 宽度
		4,   // 位数
		0.7, // 扭曲程度
		80,  // 噪点数量
	)

	captcha := base64Captcha.NewCaptcha(driver, global.CaptchaStore)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		logrus.Error(err)
		res.FailWithMsg("图片验证码生成失败", c)
		return
	}

	// 可选：打印一眼看看
	//fmt.Println("验证码 ID:", id)
	//fmt.Println("IMG:", b64s)

	res.OkWithData(CaptchaResponse{
		CaptchaID: id,
		Captcha:   b64s,
	}, c)
}
