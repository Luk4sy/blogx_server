package article_api

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"minibar_server/common/res"
	"minibar_server/global"
	"minibar_server/middleware"
	"minibar_server/models"
	"minibar_server/models/ctype"
	"minibar_server/models/enum"
	"minibar_server/utils/jwts"
	"minibar_server/utils/markdown"
)

type ArticleCreateRequest struct {
	Title       string             `json:"title"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"required,oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateRequest](c)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 判断分类 id 是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.DB.Take(&category, "id = ? and user_id = ?", *cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}

	// 文章正文防 xss 注入
	cr.Content = bluemonday.UGCPolicy().Sanitize(cr.Content)

	// 如果清洗完发现内容没了（说明用户发的全是脚本），报错
	if cr.Content == "" {
		res.FailWithMsg("文章内容包含非法字符或为空", c)
		return
	}

	// 如果未传如简介，自动从正文中取前 50 个字符
	if cr.Abstract == "" {
		// 把 markdown 转成 html，再取文本
		cr.Abstract = markdown.GetAbstract(cr.Content, 50)
	}

	// TODO:正文内容图片转存
	// 1. 图片过多？同步做，接口耗时高，异步做，

	var article = models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		UserID:      user.ID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		OpenComment: cr.OpenComment,
		CategoryID:  cr.CategoryID,
		Status:      cr.Status,
	}
	if cr.Status == enum.ArticleStatusUnderReview && global.Config.Site.Article.IsFreeReview {
		article.Status = enum.ArticleStatusPublished
	}

	err = global.DB.Create(&article).Error
	if err != nil {
		res.FailWithMsg("文章创建失败", c)
		return
	}

	res.OkWithMsg("文章创建成功", c)

}
