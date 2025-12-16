package es_service

import (
	"blogx_server/global"
	"context"
	"github.com/sirupsen/logrus"
)

func CreateIndexV2(index, mapping string) {
	if ExistsIndex(index) {
		DeleteIndex(index)
	}
	CreateIndex(index, mapping)
}

func CreateIndex(index, mapping string) {
	_, err := global.ESClient.CreateIndex(index).BodyString(mapping).Do(context.Background())
	if err != nil {
		logrus.Errorf("%s 索引创建失败 %s", index, err)
	}
	logrus.Infof("%s 创建索引成功", index)
}

// ExistsIndex 判断索引是否存在
func ExistsIndex(index string) bool {
	exist, _ := global.ESClient.IndexExists(index).Do(context.Background())
	return exist
}

func DeleteIndex(index string) {
	_, err := global.ESClient.DeleteIndex(index).Do(context.Background())
	if err != nil {
		logrus.Errorf("%s 索引删除失败 %s", index, err)
		return
	}
	logrus.Infof("%s 删除索引成功", index)
}
