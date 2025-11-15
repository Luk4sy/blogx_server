package common

import (
	"blogx_server/global"
	"fmt"
	"gorm.io/gorm"
)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"` // 前端可以覆盖
}
type Options struct {
	PageInfo     PageInfo
	Likes        []string
	Preloads     []string
	Where        *gorm.DB
	Debug        bool
	DefaultOrder string
}

func (p PageInfo) GetPage() int {
	if p.Page >= 20 || p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 100 {
		return 10
	}
	return p.Limit
}

func (p PageInfo) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func ListQuery[T any](model T, option Options) (list []T, count int, err error) {

	// 基础查询
	query := global.DB.Model(model).Where(model)

	// 日志
	if option.Debug {
		query = query.Debug()
	}

	// 排序
	if option.PageInfo.Order != "" {
		// 在外层配置了
		query = query.Order(option.PageInfo.Order)
	} else {
		if option.DefaultOrder != "" {
			query = query.Order(option.DefaultOrder)
		}
	}

	// 模糊查询
	if len(option.Likes) > 0 && option.PageInfo.Key != "" {
		likes := global.DB.Where(model)
		for _, column := range option.Likes {
			likes.Or(
				fmt.Sprintf("%s like ?", column),
				fmt.Sprintf("%%%s%%", option.PageInfo.Key))
		}
		query = query.Where(likes)
	}

	// 预加载
	for _, preload := range option.Preloads {
		query = query.Preload(preload)
	}

	// 高级查询
	if option.Where != nil {
		query = query.Where(option.Where)
	}

	// 总数查询
	var _c int64
	global.DB.Model(model).Count(&_c)
	count = int(_c)

	// 分页查询
	limit := option.PageInfo.GetLimit()
	offset := option.PageInfo.GetOffset()
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return
}
