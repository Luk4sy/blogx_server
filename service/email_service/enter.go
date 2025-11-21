package email_service

import (
	"blogx_server/global"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

// SendRegisterCode  注册账号
func SendRegisterCode(to string, code string) error {
	em := global.Config.Email
	subject := fmt.Sprintf("[%s]账号注册", em.SendNickname)
	text := fmt.Sprintf("账号注册进行中，验证码为：%s，十分钟内有效", code)

	return SendEmail(to, subject, text)
}

// SendResetPwdCode 重置密码
func SendResetPwdCode(to string, code string) error {
	em := global.Config.Email
	subject := fmt.Sprintf("[%s]密码重置", em.SendNickname)
	text := fmt.Sprintf("账号密码重置操作进行中，验证码为：%s，十分钟内有效", code)

	return SendEmail(to, subject, text)
}

func SendEmail(to, subject, text string) (err error) {
	em := global.Config.Email

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", em.SendNickname, em.SendEmail)
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(text)
	err1 := e.Send(fmt.Sprintf("%s:%d", em.Domain, em.Port), smtp.PlainAuth("", em.SendEmail, em.AuthCode, em.Domain))
	if err1 != nil && !strings.Contains(err1.Error(), "short response:") {
		return err1
	}
	return nil
}
