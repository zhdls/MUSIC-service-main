package middleware

//直接取到方法响应主体

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/pkg/logger"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 实现了双写，因此我们可以直接通过 AccessLogWriter 的 body 取到值
func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

// AccessLog 编写访问日志的中间件
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//初始化了 AccessLogWriter，将其赋予给当前的 Writer 写入流（可理解为替换原有）
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		//得到所需的日志属性
		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}
		s := "access log: method: %s, status_code: %d, " + "begin_time: %d, end_time: %d"
		global.Logger.WithFields(fields).Infof(c, s,
			c.Request.Method,
			bodyWriter.Status(),
			beginTime,
			endTime,
		)
	}
}
