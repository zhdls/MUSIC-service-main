package middleware

//将整体的限流器与对应的中间件逻辑串联

import (
	"github.com/gin-gonic/gin"
	"github.com/travel_study/blog-service/pkg/app"
	"github.com/travel_study/blog-service/pkg/errcode"
	"github.com/travel_study/blog-service/pkg/limiter"
)

// RateLimiter 入参应该为 LimiterIface 接口类型，只要符合该接口类型的具体限流器实现都可以传入并使用
func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if bucket, ok := l.GetBucket(key); ok {
			count := bucket.TakeAvailable(1) //占用存储桶中立即可用的令牌的数量，返回值为删除的令牌数
			if count == 0 {                  //返回 0，也就是已经超出配额
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
