package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// 压缩后长度与原始 Content-Length 不再一致，必须在写响应头前移除。
func (w *gzipWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipWriter) WriteHeaderNow() {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeaderNow()
}

func (w *gzipWriter) Write(data []byte) (int, error) { return w.writer.Write(data) }

func (w *gzipWriter) WriteString(s string) (int, error) { return w.writer.Write([]byte(s)) }

// Gzip 对客户端支持 gzip 的响应做压缩，跳过预检请求。
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions || !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		writer := gzip.NewWriter(c.Writer)
		defer writer.Close()

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{ResponseWriter: c.Writer, writer: writer}
		c.Next()
	}
}
