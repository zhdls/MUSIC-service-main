package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"           //多语言包 该库是与 universal-translator 配套使用的。
	"github.com/go-playground/locales/zh"           //多语言包 该库是与 universal-translator 配套使用的。
	"github.com/go-playground/locales/zh_Hant_TW"   //多语言包 该库是与 universal-translator 配套使用的。
	"github.com/go-playground/universal-translator" //通用翻译器，是一个使用 CLDR 数据 + 复数规则的 Go 语言 i18n 转换器。
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en" //validator 的翻译器
	zh_translations "github.com/go-playground/validator/v10/translations/zh" //validator 的翻译器
)

func Translations() gin.HandlerFunc { //翻译
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		//c.GetHeader获取约定的 header 参数 locale，用于判别当前请求的语言类别是 en 又或是 zh
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			}
			c.Set("trans", trans)
		}

		c.Next()
	}
}
