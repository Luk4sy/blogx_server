package flags

import (
	"blogx_server/global"
	"blogx_server/models"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.UserConfModel{},
		&models.ArticleModel{},
		&models.CategoryModel{},
		&models.ArticleDiggModel{},
		&models.CollectModel{},
		&models.UserArticleCollectModel{},
		&models.UserTopArticleModel{},
		&models.ImageModel{},
		&models.UserArticleLookHistoryModel{}, // 用户浏览文章的历史表
		&models.CommentModel{},
		&models.BannerModel{},
		&models.LogModel{}, //日志表
		&models.UserLoginModel{},
		&models.GlobalNotificationModel{},
		&models.UserLoginModel{}, // 用户登录记录表
	)
	if err != nil {
		logrus.Errorf("数据库迁移失败 %s", err)
		return
	}
	logrus.Infof("数据库迁移成功")
}
