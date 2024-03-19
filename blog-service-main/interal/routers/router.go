package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/interal/middleware"
	"github.com/travel_study/blog-service/interal/routers/api"
	"github.com/travel_study/blog-service/pkg/limiter"
	"net/http"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()
	//链路追踪   该中间件要在所有路由方法之前生效！！！！害我查错查半天www....
	r.Use(middleware.Tracing())
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	//初始化 docs 包和注册一个针对 swagger 的路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//限流
	r.Use(middleware.RateLimiter(methodLimiters))
	//对请求进行超时时间控制
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))

	//注册翻译器中间件
	r.Use(middleware.Translations())

	//新增POST方法的 /upload/file 路由	上传图片和文件服务
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	//实现 File Server , 设置文件服务去提供静态资源的访问，才能实现让外部请求本项目 HTTP Server 时同时提供静态资源的访问
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	//新增 auth 相关路由（JWT）
	r.POST("/auth", api.GetAuth)

	return r
}
