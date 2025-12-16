package models

import (
	_ "embed"
	"minibar_server/models/ctype"
)

type ArticleModel struct {
	Model
	Title        string     `gorm:"size:32" json:"title"`
	Abstract     string     `gorm:"size:256" json:"abstract"`
	Content      string     `json:"content"`
	CategoryID   uint       `json:"categoryID"`                   // 分类的id
	TagList      ctype.List `gorm:"type:longtext" json:"tagList"` // 标签列表
	Cover        string     `gorm:"size:256" json:"cover"`
	UserID       uint       `json:"userID"`
	UserModel    UserModel  `gorm:"foreignKey:UserID" json:"-"` // 关联的用户信息
	LookCount    int        `json:"lookCount"`                  // 阅读数
	DiggCount    int        `json:"diggCount"`                  // 点赞数
	CommentCount int        `json:"commentCount"`               // 评论数
	CollectCount int        `json:"collectCount"`               // 收藏数
	OpenComment  bool       `json:"openComment"`                // 开启评论
	Status       int8       `json:"status"`                     // 状态 草稿 审核中 已发布
}

//go:embed mappings/article_mapping.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}

func (ArticleModel) Index() string {
	return "article_index"
}
