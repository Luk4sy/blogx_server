package core

import (
	"blogx_server/global"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func EsConnect() *elastic.Client {
	es := global.Config.ES
	if es.Url == "" {
		return nil
	}
	client, err := elastic.NewClient(
		elastic.SetURL(es.Url),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(es.Username, es.Password),
	)
	if err != nil {
		logrus.Panicf("es 连接失败 %s", err)
		return nil
	}
	logrus.Info("es 连接成功")
	return client

}
