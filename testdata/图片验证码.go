package main

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

type CaptchaResponse struct {
	Id     string
	Encode string
}

func main() {
	resp := Captcha()
	if resp == nil {
		fmt.Println("生成验证码失败")
		return
	}

	// 输出 ID 和 base64 图片字符串
	fmt.Println("验证码 ID:", resp.Id)
	fmt.Println("IMG:", resp.Encode)

	// 复制 resp.Encode 到浏览器地址栏检查效果：
	// data:image/png;base64,你的字符串
}

func Captcha() *CaptchaResponse {
	// 配置验证码的参数
	// 使用内置的数字验证码驱动，先确认功能正常
	driver := base64Captcha.NewDriverDigit(
		60,  // 高度
		200, // 宽度
		4,   // 位数
		0.7, // 扭曲程度
		80,  // 噪点数量
	)

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil
	}
	fmt.Println("正确答案:", answer)

	return &CaptchaResponse{
		Id:     id,
		Encode: b64s,
	}
}

func VerifyCaptcha(id string, VerifyValue string) bool {
	if store.Verify(id, VerifyValue, true) {
		//验证成功
		return true
	}
	return false
}
