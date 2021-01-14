package fluentd

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
)

const (
	CONFIG_FLUENTD_HOST = "logger.fluentd.host"
	CONFIG_FLUENTD_PORT = "logger.fluentd.port"
	FLUENTD_TAG_PREFIX  = "roav_backend."
)

func BuildFluentdHook(componentName string) logrus.Hook {
	fluentdHost := beego.AppConfig.String(CONFIG_FLUENTD_HOST)
	fluentdPort, err := beego.AppConfig.Int(CONFIG_FLUENTD_PORT)
	if err != nil {
		fmt.Errorf("Get fluentdPort error ", err)
		return nil
	}
	if fluentdHook, err := New(fluentdHost, fluentdPort, FLUENTD_TAG_PREFIX+componentName); err == nil {
		return fluentdHook
	} else {
		fmt.Errorf("New fluentdHook error ", err)
		return nil
	}
}
