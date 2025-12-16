package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"minibar_server/global"
	"os"
	"path"
	"time"
)

// 颜色常量
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

// Format 实现Formatter接口 (仅用于控制台输出，带颜色)
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// 带颜色的格式
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s\n", timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

type FileDateHook struct {
	file     *os.File
	logPath  string
	fileDate string // 判断日期切换目录
	appName  string
}

func (hook *FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 注意：这里必须是指针接收器 *FileDateHook
func (hook *FileDateHook) Fire(entry *logrus.Entry) error {
	// 获取当前日期的字符串，用于判断是否需要新建文件夹
	timer := entry.Time.Format("2006-01-02")

	// 【优化】手动拼接不带颜色的日志字符串写入文件
	// 如果直接用 entry.String() 会把控制台的颜色代码也写进去，导致文件乱码
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var msg string
	if entry.HasCaller() {
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		msg = fmt.Sprintf("[%s] [%s] %s %s\n", timestamp, entry.Level, fileVal, entry.Message)
	} else {
		msg = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	// 如果日期一样，直接写入
	if hook.fileDate == timer {
		hook.file.Write([]byte(msg))
		return nil
	}

	// === 日期变了，需要切割文件 ===

	// 1. 关闭旧文件
	hook.file.Close()

	// 2. 创建新目录 (按照 logs/2025-12-12 这种格式)
	newDir := fmt.Sprintf("%s/%s", hook.logPath, timer)
	os.MkdirAll(newDir, os.ModePerm)

	// 3. 拼接新文件名
	filename := fmt.Sprintf("%s/%s.log", newDir, hook.appName)

	// 4. 打开新文件
	// O_WRONLY: 只写, O_APPEND: 追加, O_CREATE: 不存在则创建
	newFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)

	// 5. 更新 hook 的状态
	hook.file = newFile
	hook.fileDate = timer

	// 6. 写入当前这条日志
	hook.file.Write([]byte(msg))

	return nil
}

func InitFile(logPath, appName string) {
	// 初始状态，按照当前日期创建目录
	timer := time.Now().Format("2006-01-02")
	dir := fmt.Sprintf("%s/%s", logPath, timer)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return
	}

	filename := fmt.Sprintf("%s/%s.log", dir, appName)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 这里传入指针 &FileDateHook
	fileHook := &FileDateHook{file, logPath, timer, appName}
	logrus.AddHook(fileHook)
}

func InitLogrus() {
	logrus.SetOutput(os.Stdout)          // 设置输出类型为标准输出
	logrus.SetReportCaller(true)         // 开启返回函数名和行号
	logrus.SetFormatter(&LogFormatter{}) // 设置自定义的控制台 Formatter
	logrus.SetLevel(logrus.DebugLevel)   // 设置最低的 Level

	// 读取配置
	l := global.Config.Log
	InitFile(l.Dir, l.App)
}
