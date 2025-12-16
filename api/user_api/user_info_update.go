package user_api

import (
	"github.com/gin-gonic/gin"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/utils/jwts"
	"minibar_server/utils/maps"
	"time"
)

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`  // 公开我的收藏
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`    // 公开我的关注
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`        // 公开我的粉丝
	HomeStyleID *uint     `json:"homeStyleID" s-u-c:"home_style_id"` // 主页样式的id
}

func (UserApi) UserInfoUpdateView(c *gin.Context) {
	// TODO:用户头像或者名称违规处理
	var cr UserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	userMap := maps.StructToMap(cr, "s-u")
	userConfMap := maps.StructToMap(cr, "s-u-c")

	claims := jwts.GetClaims(c)

	// 【修改点 1】：把查询 userModel 移到最外面！
	// 无论是否更新 userMap，我们都需要先拿到当前用户的数据（因为后面更新配置也需要用到 UserConfModel 的 ID）
	var userModel models.UserModel
	err = global.DB.Preload("UserConfModel").Take(&userModel, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// ================= 处理 User 表更新 =================
	if len(userMap) > 0 {
		// 判断是否修改了用户名
		if cr.Username != nil && *cr.Username != userModel.Username {
			var userCount int64
			global.DB.Model(models.UserModel{}).
				Where("username = ? and id <> ?", *cr.Username, claims.UserID).
				Count(&userCount)
			if userCount > 0 {
				res.FailWithMsg("用户名被使用", c)
				return
			}

			// 检查距离上次修改是否超过 30 天 (720小时)
			var uud = userModel.UserConfModel.UpdateUsernameDate
			if uud != nil {
				if time.Since(*uud).Hours() < 720 {
					res.FailWithMsg("用户名30天内只能修改一次", c)
					return
				}
			}
			// 如果用户名改了，记录修改时间到 userConfMap 中
			// 注意：这里 userConfMap 可能会因此增加一个字段
			userConfMap["update_username_date"] = time.Now()
		}

		err = global.DB.Model(&userModel).Updates(userMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}

	// ================= 处理 UserConf 表更新 =================
	// 【修改点 2】：把这段逻辑拿出来，变成“并列”关系，而不是“嵌套”关系
	// 这样即使 userMap 为空（只改配置），这里也能执行到
	if len(userConfMap) > 0 {
		// 【修改点 3】：直接使用 userModel.UserConfModel
		// 因为上面 Preload 已经加载出来了，不需要再 global.DB.Take 查一次了
		// GORM 更新必须传指针 (&Struct)，否则不知道表名
		err = global.DB.Model(&userModel.UserConfModel).Updates(userConfMap).Error
		if err != nil {
			res.FailWithMsg("用户配置信息修改失败", c)
			return
		}
	}

	res.OkWithMsg("用户信息修改成功", c)
}
