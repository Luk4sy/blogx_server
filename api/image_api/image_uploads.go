package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

func (ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	// 文件大小判断
	s := global.Config.Upload.Size
	if fileHeader.Size > s*1024*1024 {
		res.FailWithMsg(fmt.Sprintf("文件大小大于 %dMB", s), c)
		return
	}
	// 后缀判断
	filename := fileHeader.Filename
	sufix, err := ImageSuffixJudge(filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	// 文件 Hash
	file, err := fileHeader.Open()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	byteData, _ := io.ReadAll(file)
	hash := utils.Md5(byteData)
	// 判断 Hash 有没有
	var model models.ImageModel
	err = global.DB.Take(&model, "hash = ?", hash).Error
	if err == nil {
		// 找到了
		logrus.Infof("上传图片重复 %s <==> %s  %s", filename, model.Filename, hash)
		res.Ok(model.WebPath(), "上传成功", c)
		return
	}

	// 文件名一样，但内容不一样

	filePath := fmt.Sprintf("uploads/%s/%s.%s", global.Config.Upload.UploadDir, hash, sufix)
	// 入库
	model = models.ImageModel{
		Filename: filename,
		Path:     filePath,
		Size:     fileHeader.Size,
		Hash:     hash,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	c.SaveUploadedFile(fileHeader, filePath)
	res.Ok(model.WebPath(), "图片上传成功", c)
}

func ImageSuffixJudge(filename string) (suffix string, err error) {
	_list := strings.Split(filename, ".")
	if len(_list) == 1 {
		err = errors.New("错误的文件名")
	}
	suffix = _list[len(_list)-1]
	if !utils.InList(suffix, global.Config.Upload.WhiteList) {
		err = errors.New("文件非法")
		return
	}
	return
}
