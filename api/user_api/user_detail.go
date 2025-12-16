package user_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
	"time"
)

type UserDetailResponse struct {
	ID             uint                    `json:"id"`
	CreatedAt      time.Time               `json:"createdAt"`
	Username       string                  `json:"username"`
	Nickname       string                  `json:"nickname"`
	Avatar         string                  `json:"avatar"`
	Abstract       string                  `json:"abstract"`
	RegisterSource enum.RegisterSourceType `json:"registerSource"`
	CodeAge        int                     `json:"codeAge"`
	models.UserConfModel
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var user models.UserModel
	err := global.DB.Preload("UserConfModel").Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	var data = UserDetailResponse{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Abstract:       user.Abstract,
		RegisterSource: user.RegisterSource,
		CodeAge:        user.CodeAge(),
	}
	if user.UserConfModel != nil {
		data.UserConfModel = *user.UserConfModel
	}

	res.OkWithData(data, c)
}
