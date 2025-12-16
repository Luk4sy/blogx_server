package core

import (
	"blogx_server/global"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDB() *gorm.DB {
	if len(global.Config.DB) == 0 {
		logrus.Fatalf("未配置数据库")
	}
	dc := global.Config.DB[0]

	// test
	fmt.Printf("DB config: host=%q port=%d user=%q db=%q dsn=%q\n",
		dc.Host, dc.Port, dc.User, dc.DB, dc.DSN())
	db, err := gorm.Open(mysql.Open(dc.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
	})
	if err != nil {
		logrus.Fatalf("连接数据库失败  %s", err)
	}

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")

	if len(global.Config.DB) > 1 {
		var readList []gorm.Dialector
		for _, d := range global.Config.DB[1:] {
			readList = append(readList, mysql.Open(d.DSN()))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			// use `db2` as sources, `db3`, `db4` as replicas
			Sources:  []gorm.Dialector{mysql.Open(dc.DSN())}, //写
			Replicas: readList,                               //读
			// sources/replicas load balancing policy
			Policy: dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logrus.Fatalf("读写配置错误 %s", err)
		}
	}
	return db
}
