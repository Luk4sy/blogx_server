package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"strings"
)

type ActionLog struct {
	c            *gin.Context
	level        enum.LogLevelType
	title        string
	requestBody  []byte
	responseBody []byte
	log          *models.LogModel
	showRequest  bool
	showResponse bool
	itemList     []string
}

func (ac *ActionLog) ShowRequest() {
	ac.showResponse = true
}

func (ac *ActionLog) ShowResponse() {
	ac.showRequest = true
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}
func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

func (ac *ActionLog) setItem(label string, value any, levelType enum.LogLevelType) {

	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, _ := json.Marshal(value)
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	ac.itemList = append(ac.itemList, fmt.Sprintf(`<div class="log_item %s"><div class="log_item_label">%s</div><div class="log_item_content">%s</div></div>`,
		levelType,
		label,
		v,
	))
}

func (ac *ActionLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}

func (ac *ActionLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)
}

func (ac *ActionLog) SetRequest(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	ac.requestBody = byteData
}

func (ac *ActionLog) SetResponse(data []byte) {
	ac.responseBody = data
}

func (ac *ActionLog) Save() {

	if ac.log != nil {
		// 之前已经save过了，那就是更新
		global.DB.Model(ac.log).Updates(map[string]any{
			"title": "更新",
		})
		return
	}

	var newItemList []string
	// 设置请求
	if ac.showRequest {
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_request\">\n<div class=\"log_request_head\">\n<span class=\"log_request_method delete\">%s</span>\n<span class=\"log_request_path\">/%s</span>\n</div>\n<div class=\"log_request_body\">\n<pre class=\"log_json_body\">%s</pre>\n</div>\n</div>\n", ac.c.Request.Method, ac.c.Request.URL.String(), string(ac.requestBody)))
	}

	// 中间一些content
	newItemList = append(newItemList, ac.itemList...)

	// 设置响应
	if ac.showResponse {
		newItemList = append(newItemList, fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>\n",
			string(ac.responseBody)))
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	userID := uint(1)

	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(ac.itemList, "\n"),
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}

	err := global.DB.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err)
		return
	}
	ac.log = &log

}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}

func GetLog(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		return NewActionLogByGin(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLogByGin(c)
	}
	return log
}
