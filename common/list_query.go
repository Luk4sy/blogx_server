package common

import (
	"fmt"
	"gorm.io/gorm"
	"minibar_server/global"
)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"` // 前端可以覆盖
}

// Options 查询规则
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

// ListQuery 封装排序、模糊搜索、预加载、额外 where、分页、统计总数
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
		likes := global.DB.Session(&gorm.Session{NewDB: true}) // 新起一个干净 DB
		key := fmt.Sprintf("%%%s%%", option.PageInfo.Key)

		for i, column := range option.Likes {
			cond := fmt.Sprintf("%s LIKE ?", column)
			if i == 0 {
				likes = likes.Where(cond, key)
			} else {
				likes = likes.Or(cond, key)
			}
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
	if err = query.Count(&_c).Error; err != nil {
		return nil, 0, err
	}
	count = int(_c)

	// 分页查询
	limit := option.PageInfo.GetLimit()
	offset := option.PageInfo.GetOffset()
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return
}
