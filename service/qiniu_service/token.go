package qiniu_service

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"minibar_server/global"
	"time"
)

func GenToken() (token string, err error) {
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)
	putPolicy, err := uptoken.NewPutPolicy(global.Config.QiNiu.Bucket, time.Now().Add(time.Duration(global.Config.QiNiu.Expiry)*time.Second))
	if err != nil {
		return
	}
	token, err = uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return
	}
	return
}
