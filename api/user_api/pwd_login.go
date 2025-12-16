package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/service/user_service"
	"minibar_server/utils/jwts"
	"minibar_server/utils/pwd"
)

type PwdLoginRequest struct {
	Val      string `json:"val"  binding:"required"` // 可能是用户名 也可能是邮箱
	Password string `json:"password" binding:"required"`
}

func (UserApi) PwdLoginApi(c *gin.Context) {

	cr := middleware.GetBind[PwdLoginRequest](c)
	fmt.Println(cr)

	if !global.Config.Site.Login.UsernamePwdLogin {
		res.FailWithMsg("站点未启用密码登录", c)
		return
	}

	var user models.UserModel
	err := global.DB.Take(&user, "(username = ? or email = ?) and password <> ''",
		cr.Val, cr.Val).Error
	if err != nil {
		res.FailWithMsg("用户名密码错误", c)
		return
	}
	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		res.FailWithMsg("用户名密码错误", c)
		return
	}

	// 颁发 token
	token, _ := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	user_service.NewUserService(user).UserLogin(c)

	res.OkWithData(token, c)
}
