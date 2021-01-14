package localUtil

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

var curLang string = ""

func GetString(key string, args ...interface{}) string {
	if curLang == "" {
		curLang = beego.AppConfig.String("lang")
	}
	return i18n.Tr(curLang, key, args...)
}

func Tr(lang, key string, args ...interface{}) string {
	return i18n.Tr(lang, key, args...)
}

//Get Locale Strings with Arguments
func GetStringWithArgs(key string, args ...interface{}) string {
	if curLang == "" {
		curLang = beego.AppConfig.String("lang")
	}
	return i18n.Tr(curLang, key, args)
}
