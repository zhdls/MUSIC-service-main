package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

// NewSetting 初始化本项目的配置的基础属性
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	//vp.AddConfigPath("configs/")
	//vp.AddConfigPath("../blog-service/configs")
	//vp.AddConfigPath("./configs")//根据当前终端执行目录的相对位置
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	//新增文件热更新的监听和变更处理
	s := &Setting{vp}
	s.WatchSettingChange()

	return s, nil
}

func (s *Setting) WatchSettingChange() {
	//起一个协程
	go func() {
		s.vp.WatchConfig()                            //对文件配置进行监听
		s.vp.OnConfigChange(func(in fsnotify.Event) { //配置文件变化了
			_ = s.ReloadAllSection() //把最新的配置信息反序列化到Conf中
		})
	}()
}
