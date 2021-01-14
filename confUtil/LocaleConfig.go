package confUtil

import (
	//"yproject/github.com/yiGmMk/pz-infra-new/log"

	//"github.com/astaxie/beego"
	//"github.com/beego/i18n"
	//"github.com/go-sql-driver/mysql"
	"github.com/yiGmMk/pz-infra-new/log"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func LoadLocaleConfig() {
	lang := beego.AppConfig.String("lang")
	if lang == "" {
		log.Info("no lang in configuraiton, skip locale")
		return
	}
	configFileName := "conf/locale_" + lang + ".ini"
	localeFilePath := FileHierarchyFind(configFileName)
	log.Info("localeFilePath lang in configuraiton, " + localeFilePath)
	if localeFilePath == "" {
		log.Warnf("not found Locale in current or parent directories, configFileName %s", configFileName)
		return
	}
	if err := i18n.SetMessage(lang, localeFilePath); err != nil {
		log.Errorf("Failed to load Locale , localeFilePath %s, %+v", localeFilePath, err)
	}
	log.Debugf("load localization ini file, localeFilePath %s", localeFilePath)
}
