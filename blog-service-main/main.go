package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/interal/model"
	"github.com/travel_study/blog-service/interal/routers"
	"github.com/travel_study/blog-service/pkg/logger"
	"github.com/travel_study/blog-service/pkg/setting"
	"github.com/travel_study/blog-service/pkg/tracer"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//程序执行顺序是：全局变量初始化->init方法->main 方法

//router -> v1 ->service -> dao ->model

var (
	port    string
	runMode string
	config  string

	isVersion bool

	buildTime    string
	buildVersion string
	gitCommitID  string
)

// 在main方法之前自动执行
// 把配置文件内容映射到应用配置结构体
func init() {
	//setupFlag一定要写在setupSetting前，才能在newSetting中顺利传入命令行参数
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	//日志组件流程
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	//链路追踪
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}

}

// @title 音乐软件
func main() {
	//设置gin的运行模式
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	//通过自定义 http.Server，设置了监听的 TCP Endpoint、处理的程序、允许读取/写入的最大时间、请求头的最大字节数等基础参数
	//设置成yaml文件中已经映射到结构体中的字段值
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err:%v", err)
		}
	}()
	//等待中断信号
	quit := make(chan os.Signal)
	//quit接收 syscall.SIGINT 和 syscall.SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	<-quit
	log.Println("Shutting down server...")

	//最大时间控制，通知该服务端 它有5s的时间来处理原有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //创建一个5秒超时的context
	defer cancel()
	//5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:%v", err)
	}

	log.Println("Server exiting")
}

func setupDBEngine() error {
	var err error
	//注意不要用:=，否则会新生成一个global.DBEngine；不会修改全局变量里的global.DBEngine，它仍然是 nil
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupSetting() error {
	setting, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	//配置项的读取和映射
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.JWTSetting.Expire *= time.Second

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	//fmt.Printf("%v\n",global.ServerSetting)
	//fmt.Printf("%v\n",global.AppSetting)
	//fmt.Printf("%v\n",global.DatabaseSetting)

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	return nil
}

// 对global的包全局变量Logger进行了初始化
func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{ //使用lumberjack作为日志库的io.Writer
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,  //设置日志文件所允许的最大占用空间为 600MB
		MaxAge:    10,   //日志文件最大生存周期为 10 天
		LocalTime: true, //设置日志文件名的时间格式为本地时间
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

// jaeger分布式链路追踪系统
func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	//flag.StringVar(&config,"config","./configs","指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
