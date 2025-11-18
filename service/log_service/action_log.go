package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwts"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// ActionLog：用于记录“接口级操作日志”的结构体
// 典型场景：某个 API 想记录请求/响应 + 额外信息（结构体、错误等），最终写入 log_models 表。

type ActionLog struct {
	// Gin 上下文，用来取 Request、IP、Header、Token 等
	c *gin.Context

	// 日志基础信息
	level enum.LogLevelType
	title string

	// 请求 / 响应原始数据
	requestBody    []byte
	responseBody   []byte
	responseHeader http.Header

	// 数据库里的那条日志记录（第一次 Save 之后会有值）
	log *models.LogModel

	// 开关类字段：控制是否展示哪些信息
	showRequestHeader  bool
	showResponseHeader bool
	showRequest        bool
	showResponse       bool

	// 中间累积的 HTML 片段
	itemList []string

	// 是否是“中间件触发创建”的标记
	// 用来区分：是在视图里 Save，还是在中间件里 Save。
	isMiddleware bool
}

//////////////////////////
// 一、基础配置：标题/等级 & 显示哪些内容
//////////////////////////

// ShowRequest 开启“在日志中展示请求体”
func (ac *ActionLog) ShowRequest() {
	ac.showRequest = true
}

// ShowResponse 开启“在日志中展示响应体”
func (ac *ActionLog) ShowResponse() {
	ac.showResponse = true
}

// ShowRequestHeader 开启“展示请求头”
func (ac *ActionLog) ShowRequestHeader() {
	ac.showRequestHeader = true
}

// ShowResponseHeader 开启“展示响应头”
func (ac *ActionLog) ShowResponseHeader() {
	ac.showResponseHeader = true
}

// SetTitle 设置日志标题（对应 log_models.title）
func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

// SetLevel 设置日志级别（info / warn / error）
func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

//////////////////////////
// 二、附加信息：普通字段 / 链接 / 图片 / 错误
//////////////////////////

// setItem 内部方法：把一个 label+value 以 HTML 形式追加到 itemList
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

	ac.itemList = append(ac.itemList, fmt.Sprintf(
		`<div class="log_item %s"><div class="log_item_label">%s</div><div class="log_item_content">%s</div></div>`,
		levelType,
		label,
		v,
	))
}

// SetItem 是 SetItemInfo 的别名，默认 info 级别
func (ac *ActionLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

// SetItemInfo 以 info 级别追加一条普通信息
func (ac *ActionLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

// SetItemWarn 以 warn 级别追加一条信息
func (ac *ActionLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}

// SetItemError 以 error 级别追加一条信息
func (ac *ActionLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)
}

// SetLink 追加一个可点击的超链接
func (ac *ActionLog) SetLink(label string, href string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf(
		`<div class="log_item link"><div class="log_item_label">%s</div><div class="log_item_content"><a href="%s" target="_blank">%s</a></div></div>`,
		label,
		href,
		href,
	))
}

// SetImage 追加一张图片
func (ac *ActionLog) SetImage(src string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf(
		`<div class="log_image"><img src="%s" alt=""></div>`,
		src,
	))
}

// SetError 追加一条带堆栈的错误信息：既向 logrus 输出 error，又把堆栈写入 HTML
func (ac *ActionLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s %s", label, err.Error())

	ac.itemList = append(ac.itemList, fmt.Sprintf(
		`<div class="log_error"><div class="line"><div class="label">%s</div><div class="value">%s</div><div class="type">%T</div></div><div class="stack">%+v</div></div>`,
		label,
		err,
		err,
		msg,
	))
}

//////////////////////////
// 三、请求 / 响应数据采集（给中间件调用）
//////////////////////////

// SetRequest 从 Gin 的 Request.Body 中读取一份副本，避免后续绑定出问题
func (ac *ActionLog) SetRequest(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	// 读完之后，要把 Body 再塞回去，否则后面 c.ShouldBindXXX 会读不到
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

	ac.requestBody = byteData
}

// SetResponse 在中间件中调用：记录响应 body
func (ac *ActionLog) SetResponse(data []byte) {
	ac.responseBody = data
}

// SetResponseHeader 在中间件中调用：记录响应 header
func (ac *ActionLog) SetResponseHeader(header http.Header) {
	ac.responseHeader = header
}

//////////////////////////
// 四、中间件收尾逻辑：MiddlewareSave
//////////////////////////

// MiddlewareSave 只在日志中间件的“后半段”调用
// 作用：判断这次请求是否需要落库；如果需要，则补全响应信息并调用 Save。
func (ac *ActionLog) MiddlewareSave() {
	// 看看视图层有没有调用 GetLog，把 saveLog 打成 true
	_saveLog, _ := ac.c.Get("saveLog")
	saveLog, _ := _saveLog.(bool)
	if !saveLog {
		// 视图层没有显式使用日志系统，这次请求就不落库
		return
	}

	// ac.log == nil 表示还没有创建过那条日志记录
	if ac.log == nil {
		// 由中间件触发的首次 Save，用 isMiddleware 标记一下
		ac.isMiddleware = true
		ac.Save()
		return
	}

	// 走到这里说明：
	// 视图层已经调用过 Save() 创建过日志了，这里属于“追加响应信息”的更新流程。
	if ac.showResponseHeader {
		byteData, _ := json.Marshal(ac.responseHeader)
		ac.itemList = append(ac.itemList, fmt.Sprintf(
			`<div class="log_response_header"><pre class="log_json_body">%s</pre></div>`,
			string(byteData),
		))
	}

	if ac.showResponse {
		ac.itemList = append(ac.itemList, fmt.Sprintf(
			`<div class="log_response"><pre class="log_json_body">%s</pre></div>`,
			string(ac.responseBody),
		))
	}

	ac.Save()
}

//////////////////////////
// 五、真正的持久化逻辑：Save
//////////////////////////

// Save 负责“新建或更新一条日志”
// - 如果 ac.log 为 nil：创建一条新的 log_models 记录
// - 如果 ac.log 不为 nil：在原有 content 后面追加新的 itemList 内容
func (ac *ActionLog) Save() (id uint) {
	// 说明：
	// 1）有些场景会在视图里手动 Save（只写请求/中间内容）；
	// 2）中间件里再补充响应信息；
	// 3）通过 ac.log 是否为空 + isMiddleware 来区分“创建”还是“更新”。

	// 情况一：已经有日志记录了，属于“在原有内容基础上追加”
	if ac.log != nil {
		newContent := strings.Join(ac.itemList, "\n")
		content := ac.log.Content + "\n" + newContent

		global.DB.Model(ac.log).Updates(map[string]any{
			"content": content,
		})

		// 这次已经写入完毕，清空本次的临时 itemList
		ac.itemList = []string{}
		return ac.log.ID
	}

	// 情况二：第一次写入日志，准备拼一整块内容
	var newItemList []string

	// 1. 请求头
	if ac.showRequestHeader {
		byteData, _ := json.Marshal(ac.c.Request.Header)
		newItemList = append(newItemList, fmt.Sprintf(
			`<div class="log_request_header"><pre class="log_json_body">%s</pre></div>`,
			string(byteData),
		))
	}

	// 2. 请求体
	if ac.showRequest {
		newItemList = append(newItemList, fmt.Sprintf(
			`<div class="log_request">
						<div class="log_request_head">
							<span class="log_request_method %s">%s</span>
							<span class="log_request_path">/%s</span>
						</div>
						<div class="log_request_body">
						<pre class="log_json_body">%s</pre>
						</div>
					</div>
`,
			strings.ToLower(ac.c.Request.Method),
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.requestBody),
		))
	}

	// 3. 中间追加的内容（结构体、错误、链接、图片等）
	newItemList = append(newItemList, ac.itemList...)

	// 4. 如果是“中间件触发创建”，需要在这里同时拼上响应头和响应体
	if ac.isMiddleware {
		if ac.showResponseHeader {
			byteData, _ := json.Marshal(ac.responseHeader)
			newItemList = append(newItemList, fmt.Sprintf(
				`<div class="log_response_header"><pre class="log_json_body">%s</pre></div>`,
				string(byteData),
			))
		}

		if ac.showResponse {
			newItemList = append(newItemList, fmt.Sprintf(
				`<div class="log_response"><pre class="log_json_body">%s</pre></div>`,
				string(ac.responseBody),
			))
		}
	}

	// 5. 补全 IP 和用户信息
	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)

	claims, err := jwts.ParseTokenByGin(ac.c)
	userID := uint(0)
	if err == nil && claims != nil {
		userID = claims.UserID
	}

	// 6. 真正落库
	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(newItemList, "\n"),
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}

	if err = global.DB.Create(&log).Error; err != nil {
		logrus.Errorf("日志创建失败 %s", err)
		return
	}

	ac.log = &log
	ac.itemList = []string{}
	return log.ID
}

//////////////////////////
// 六、构造 & 获取：给中间件 / 视图层使用
//////////////////////////

// NewActionLogByGin 创建一个绑定当前 gin.Context 的 ActionLog
func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}

// GetLog 从 gin.Context 中取出日志对象，如果没有则新建一个。
// 同时会把 "saveLog" 标记成 true，表示这次请求需要落库。
func GetLog(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		return NewActionLogByGin(c)
	}

	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLogByGin(c)
	}

	// 标记这次请求需要保存日志，中间件的 MiddlewareSave 会用到
	c.Set("saveLog", true)
	return log
}
