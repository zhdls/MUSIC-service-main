package middleware

//统一的在应用程序中针对所有请求都进行一个最基本的超时时间控制

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

// ContextTimeout 上下文超时时间控制的中间件
func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t) //设置当前 context 的超时时间,并重新赋给了 gin.Context
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
