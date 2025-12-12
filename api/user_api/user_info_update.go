package user_api

import (
	"blogx_server/common/res"
	"github.com/gin-gonic/gin"
)

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username"`
	Nickname    *string   `json:"nickname"`
	Avatar      *string   `json:"avatar"`
	Abstract    *string   `json:"abstract"`
	LikeTags    *[]string `json:"likeTags"`
	OpenCollect *bool     `json:"openCollect"` // 公开我的收藏
	OpenFollow  *bool     `json:"openFollow"`  // 公开我的关注
	OpenFans    *bool     `json:"openFans"`    // 公开我的粉丝
	HomeStyleID *uint     `json:"homeStyleID"` // 主页样式的id
}

func (UserApi) UserInfoUpdateView(c *gin.Context) {
	var cr UserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
}
