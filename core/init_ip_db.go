package core

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
	ipUtils "minibar_server/utils/ip"
	"strings"
)

var searcher *xdb.Searcher

func InitIPDB() {
	var dbPath = "init/ip2region.xdb"
	_searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		logrus.Fatalf("ip地址数据库加载失败 %s", err)
		return
	}
	searcher = _searcher
}

func GetIpAddr(ip string) (addr string) {
	if ipUtils.HasLocalIPAddr(ip) {
		return "此 ip 为内网 ip"
	}
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logrus.Warnf("错误的ip %s", err)
		return "异常地址"
	}
	_addList := strings.Split(region, "|")
	if len(_addList) != 5 {
		logrus.Warnf("异常的 ip 地址 %s", ip)
		return "未知地址"
	}

	// _addrList 五个部分判断
	// 国家 0 省份 城市 运营商
	country := _addList[0]
	province := _addList[2]
	city := _addList[3]

	if province != "0" && city != "0" {
		return fmt.Sprintf("%s·%s", province, city)
	}
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s·%s", country, province)
	}
	if country != "0" {
		return country
	}

	return region
}
