package middleware

import (
	"github.com/gin-gonic/gin"
	"minibar_server/service/log_service"
	"net/http"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Body []byte
	Head http.Header
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func (w *ResponseWriter) Header() http.Header {
	return w.Head
}

func LogMiddleware(c *gin.Context) {
	log := log_service.NewActionLogByGin(c)

	// ① 进来时：记录请求体
	log.SetRequest(c)
	c.Set("log", log)

	// 包装响应 writer，方便抓响应 body 和 header
	res := &ResponseWriter{
		ResponseWriter: c.Writer,
		Head:           make(http.Header),
	}
	c.Writer = res

	// ② 放行，进入路由 + 具体 handler（也就是你的 api 层）
	c.Next()

	// ③ 出来时：记录响应
	log.SetResponse(res.Body)
	log.SetResponseHeader(res.Head)
	log.MiddlewareSave()
}
