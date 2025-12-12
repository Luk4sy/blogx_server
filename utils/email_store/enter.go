package email_store

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type EmailStoreInfo struct {
	Email     string
	Code      string
	Timestamp int64 // 添加时间戳字段
}

var emailVerifyStore = sync.Map{}

func Set(id, email, code string) {
	emailVerifyStore.Store(id, EmailStoreInfo{
		Email:     email,
		Code:      code,
		Timestamp: time.Now().Unix(), // 存储当前时间戳
	})
	logrus.Infof("存入邮箱验证码，ID: %s, Code: %s", id, code)
}

func Verify(id, code string) (info EmailStoreInfo, ok bool) {
	value, ok := emailVerifyStore.Load(id)
	if !ok {
		logrus.Infof("未找到邮箱验证码数据，ID: %s", id) // 调试日志
		return
	}
	info, ok = value.(EmailStoreInfo)
	if !ok {
		logrus.Infof("转换错误，ID: %s, value: %+v", id, value) // 调试日志
		return
	}

	logrus.Infof("从存储中读取的值: %+v", info) // 调试日志

	// 检查是否过期（有效期设定为10分钟）
	if time.Now().Unix()-info.Timestamp > 600 { // 600秒 = 10分钟
		emailVerifyStore.Delete(id)
		ok = false
		return
	}

	if info.Code != code {
		emailVerifyStore.Delete(id)
		ok = false
		return
	}

	emailVerifyStore.Delete(id)
	return
}
