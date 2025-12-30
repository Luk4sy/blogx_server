package markdown

import (
	"bytes"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MdToHtml 将 Markdown 转换为 HTML
func MdToHtml(source string) string {
	// 1. 配置 goldmark
	md := goldmark.New(
		// 支持 GFM (GitHub Flavored Markdown)：表格、删除线等
		goldmark.WithExtensions(
			extension.GFM,
			// 代码高亮扩展 (使用 dracula 风格，你可以换成 monokai, github 等)
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
		// 解析器选项
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 自动为标题生成 ID，方便做目录 TOC
		),
		// 渲染器选项
		goldmark.WithRendererOptions(
			// 允许 HTML 标签（如果你信任输入的内容，可以开启；否则建议关闭或配合 bluemonday 使用）
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	// 2. 执行转换
	if err := md.Convert([]byte(source), &buf); err != nil {
		return ""
	}

	return buf.String()
}

// GetAbstract 获取文章摘要 (Markdown -> HTML -> 纯文本 -> 截取)
// content: Markdown 原文
// length: 想要截取的长度（中文按 1 个字算）
func GetAbstract(content string, length int) string {
	// 1. 先转成 HTML
	htmlStr := MdToHtml(content)

	// 2. 剔除 HTML 标签，只留纯文本
	// bluemonday.StrictPolicy() 会把所有 <tag> 都去掉
	text := bluemonday.StrictPolicy().Sanitize(htmlStr)

	// 3. 截取字符串 (处理中文乱码问题)
	// 在 Go 中直接 text[:100] 是按字节截取，中文会乱码
	// 必须转成 []rune 切片再截取
	runes := []rune(text)

	if len(runes) > length {
		return string(runes[:length]) + "..." // 超过长度加省略号
	}

	return string(runes)
}
