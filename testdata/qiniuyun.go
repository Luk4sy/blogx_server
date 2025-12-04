package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	file2 "blogx_server/utils/file"
	"blogx_server/utils/hash"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"io"
	"time"
)

func SendFile(file string) (url string, err error) {

	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)

	hashString, err := hash.FileMd5(file)
	if err != nil {
		return
	}

	suffix, _ := file2.ImageSuffixJudge(file)
	fileName := fmt.Sprintf("%s.%s", hashString, suffix)
	key := fmt.Sprintf("%s/%s", global.Config.QiNiu.Prefix, fileName)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err = uploadManager.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: global.Config.QiNiu.Bucket,
		ObjectName: &key,
		FileName:   fileName,
	}, nil)
	return fmt.Sprintf("%s/%s", global.Config.QiNiu.Uri, key), err
}

func SendReader(reader io.Reader) (url string, err error) {

	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)

	uid := uuid.New().String()
	fileName := fmt.Sprintf("%s.png", uid)
	key := fmt.Sprintf("%s/%s", global.Config.QiNiu.Prefix, fileName)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err = uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: global.Config.QiNiu.Bucket,
		ObjectName: &key,
		FileName:   fileName,
	}, nil)
	return fmt.Sprintf("%s/%s", global.Config.QiNiu.Uri, key), err
}

func GenToken1() (token string, err error) {
	mac := credentials.NewCredentials(global.Config.QiNiu.AccessKey, global.Config.QiNiu.SecretKey)
	putPolicy, err := uptoken.NewPutPolicy(global.Config.QiNiu.Bucket, time.Now().Add(1*time.Minute))
	if err != nil {
		return
	}
	token, err = uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return
	}
	return
}

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	//url, err := SendFile("uploads/images/0ee29a20ccdde43048fdc1c7a10f874-transformed.jpeg")
	//fmt.Println(url, err)
	//file, _ := os.Open("uploads/images/0ee29a20ccdde43048fdc1c7a10f874-transformed.jpeg")
	//url, err := SendReader(file)
	//fmt.Println(url, err)

	fmt.Println(GenToken1())
}
