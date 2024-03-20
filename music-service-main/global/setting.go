package global

//全局变量

import (
	"github.com/travel_study/blog-service/pkg/logger"
	"github.com/travel_study/blog-service/pkg/setting"
)

// 针对最初预估的三个区段配置，进行了全局变量的声明，
// 便于在接下来的步骤将其关联起来，并且提供给应用程序内部调用。
var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS

	JWTSetting *setting.JWTSettingS //新增 JWT 配置的全局对象

	EmailSetting *setting.EmailSettingS // Email 对应的配置全局对象
)

var (
	Logger *logger.Logger //用于日志组件的初始化
)
