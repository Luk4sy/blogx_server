package image_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"minibar_server/common"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/service/log_service"
)

type ImageApi struct {
}

type ImageListResponse struct {
	models.ImageModel
	WebPath string `json:"webPath"`
}

// ImageListView 查询图片
func (ImageApi) ImageListView(c *gin.Context) {
	var cr common.PageInfo
	c.ShouldBindQuery(&cr)
	_list, count, _ := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"filename"},
	})

	// list 构造返回给前端的结构体
	var list = make([]ImageListResponse, 0)
	for _, model := range _list {
		list = append(list, ImageListResponse{
			ImageModel: model,
			WebPath:    model.WebPath(),
		})
	}
	res.OkWithList(list, count, c)
}

// ImageRemoveView 删除图片
func (ImageApi) ImageRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()

	var list []models.ImageModel
	global.DB.Find(&list, "id in ?", cr.IDList)

	var successCount, errCount int64
	if len(list) > 0 {
		// 删除对应的文件
		successCount = global.DB.Delete(&list).RowsAffected
		if err != nil {
			logrus.Errorf("删除图片失败 %s", err)
		}
	}

	errCount = int64(len(list)) - successCount

	msg := fmt.Sprintf("操作成功，成功删除%d条 失败%d", successCount, errCount)

	res.OkWithMsg(msg, c)
}
