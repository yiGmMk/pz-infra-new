package elastic

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"github.com/sohlich/elogrus"
	"gopkg.in/olivere/elastic.v5"
)

const (
	CONFIG_ELASTIC_ENABLE = "logger.elastic.enable"
	CONFIG_ELASTIC_URL    = "logger.elastic.url"
	CONFIG_ELASTIC_INDEX  = "logger.elastic.index"
	CONFIG_APP_HOST_NAME  = "app.host.name "
)

// Currently we are not directly forwarding logs to elastic search, but use fluentd instead
func BuildElasticHook() logrus.Hook {
	enableElastic, err := beego.AppConfig.Bool(CONFIG_ELASTIC_ENABLE)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !enableElastic {
		return nil
	}
	elasticUrl := beego.AppConfig.String(CONFIG_ELASTIC_URL)
	client, err := elastic.NewClient(elastic.SetURL(elasticUrl))
	if err != nil {
		panic(err)
	}
	elasticIndex := beego.AppConfig.String(CONFIG_ELASTIC_INDEX)
	appHostName := beego.AppConfig.String(CONFIG_APP_HOST_NAME)
	elasticHook, err := elogrus.NewElasticHook(client, appHostName, logrus.DebugLevel, elasticIndex)
	if err != nil {
		panic(err)
	}
	return elasticHook
}
