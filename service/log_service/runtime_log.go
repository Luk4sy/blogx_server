package log_service

import (
	"encoding/json"
	"fmt"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"minibar_server/global"
	"minibar_server/models"
	"minibar_server/models/enum"
	"reflect"
	"strings"
	"time"
)

type RuntimeLog struct {
	level           enum.LogLevelType
	title           string
	itemList        []string
	serviceName     string
	runtimeDateType RuntimeDateType
}

func (r *RuntimeLog) SetTitle(title string) {
	r.title = title
}
func (r *RuntimeLog) SetLevel(level enum.LogLevelType) {
	r.level = level
}

func (r *RuntimeLog) setItem(label string, value any, levelType enum.LogLevelType) {

	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, _ := json.Marshal(value)
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	r.itemList = append(r.itemList, fmt.Sprintf(`<div class="log_item %s"><div class="log_item_label">%s</div><div class="log_item_content">%s</div></div>`,
		levelType,
		label,
		v,
	))
}

func (r *RuntimeLog) SetLink(label string, href string) {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a></div></div>",
		label,
		href, href,
	))
}

func (r *RuntimeLog) SetImage(src string) {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\"></div>",
		src,
	))
}
func (r *RuntimeLog) SetItem(label string, value any) {
	r.setItem(label, value, enum.LogInfoLevel)
}

func (r *RuntimeLog) SetItemInfo(label string, value any) {
	r.setItem(label, value, enum.LogInfoLevel)
}

func (r *RuntimeLog) SetItemWarn(label string, value any) {
	r.setItem(label, value, enum.LogWarnLevel)
}

func (r *RuntimeLog) SetItemError(label string, value any) {
	r.setItem(label, value, enum.LogErrLevel)
}

func (r *RuntimeLog) SetNowTime() {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_time\">%s</div>", time.Now().Format("2006-01-02 15:04:05")))
}
func (r *RuntimeLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s %s", label, err.Error())
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"strk\">%+v</div></div>",
		label,
		err,
		err,
		msg,
	))
}

type RuntimeDateType int8

const (
	RuntimeDateHour  RuntimeDateType = 1
	RuntimeDateDay   RuntimeDateType = 2
	RuntimeDateWeek  RuntimeDateType = 3
	RuntimeDateMonth RuntimeDateType = 4
)

func (r *RuntimeLog) Save() {

	r.SetNowTime()
	// 判断创建还是更新
	var log models.LogModel
	global.DB.Find(&log, fmt.Sprintf("service_name = ? and log_type = ? and created_at >= date_sub(now(), %s)",
		r.runtimeDateType.GetSqlTime()),
		r.serviceName,
		enum.RuntimeLogType,
	)
	content := strings.Join(r.itemList, "\n")

	if log.ID != 0 {
		// 更新
		c := strings.Join(r.itemList, "\n")
		newContent := log.Content + "\n" + c

		// 之前已经save过了，那就是更新
		global.DB.Model(&log).Updates(map[string]any{
			"content": newContent,
		})
		r.itemList = []string{}
		return
	}
	err := global.DB.Create(&models.LogModel{
		LogType:     enum.RuntimeLogType,
		Title:       r.title,
		Content:     content,
		Level:       r.level,
		ServiceName: r.serviceName,
	}).Error
	if err != nil {
		logrus.Errorf("创建运行日志错误 %s", err)
		return
	}
	r.itemList = []string{}

}

func (r RuntimeDateType) GetSqlTime() string {
	switch r {
	case RuntimeDateHour:
		return "INTERVAL 1 HOUR"
	case RuntimeDateDay:
		return "INTERVAL 1 DAY"
	case RuntimeDateWeek:
		return "INTERVAL 1 WEEK"
	case RuntimeDateMonth:
		return "INTERVAL 1 MONTH"
	}
	return "INTERVAL 1 DAY"
}

func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *RuntimeLog {
	return &RuntimeLog{
		serviceName:     serviceName,
		runtimeDateType: dateType,
	}
}
