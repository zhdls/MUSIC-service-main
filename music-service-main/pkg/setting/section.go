package setting

//声明配置属性的结构体,并编写读取区段配置的配置方法

import "time"

var sections = make(map[string]interface{})

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	UploadSavePath        string
	UploadServerUrl       string
	UploadImageMaxSize    int
	UploadImageAllowExts  []string
	DefaultContextTimeout time.Duration
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

// ReadSection 读取区段配置的配置方法
func (s *Setting) ReadSection(k string, v interface{}) error {
	//读取配置文件，把key为k的配置读到v中
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	//增加读取section的存储记录，以便在重新加载配置的方法中进行二次处理
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// ReloadAllSection 重新读取配置
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
