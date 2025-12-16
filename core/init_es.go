package core

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"minibar_server/global"
)

func EsConnect() *elastic.Client {
	es := global.Config.ES

	if es.Addr == "" {
		return nil
	}
	client, err := elastic.NewClient(
		elastic.SetURL(es.EsUrl()),
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
