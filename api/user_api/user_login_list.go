package user_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/utils/jwts"
	"time"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `form:"userId"`
	Ip        string `form:"ip"`
	Addr      string `form:"addr"`
	StartTime string `form:"startTime"` // 起始时间 时分秒
	EndTime   string `form:"endTime"`
	Type      int8   `json:"type" form:"type" binding:"required,oneof=1 2"` // 1 用户：超看自己的内容 2 管理员可以查全部
}

type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickName string `json:"userNickName,omitempty"`
	UserAvatar   string `json:"userAvatar,omitempty"`
}

func (UserApi) UserLoginListView(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	claim := jwts.GetClaims(c)
	if cr.Type == 1 {
		cr.UserID = claim.UserID
	}

	var query = global.DB.Where("")
	if cr.StartTime != "" {
		_, err := time.Parse("2006-01-02 15:04:05", cr.StartTime)
		if err != nil {
			res.FailWithMsg("开始时间按格式错误", c)
			return
		}
		query.Where("created_at >= ?", cr.StartTime)
	}
	if cr.EndTime != "" {
		_, err := time.Parse("2006-01-02 15:04:05", cr.EndTime)
		if err != nil {
			res.FailWithMsg("结束时间按格式错误", c)
			return
		}
		query.Where("created_at <= ?", cr.EndTime)
	}
	var preloads []string
	if cr.Type == 2 {
		preloads = []string{"UserModel"}
	}

	_list, count, _ := common.ListQuery(models.UserLoginModel{
		UserID: cr.UserID,
		IP:     cr.Ip,
		Addr:   cr.Addr,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: preloads,
	})

	var list = make([]UserLoginListResponse, 0)
	for _, model := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: model,
			UserNickName:   model.UserModel.Nickname,
			UserAvatar:     model.UserModel.Avatar,
		})
	}

	res.OkWithList(list, count, c)

}
