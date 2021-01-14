package log

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/cihub/seelog"
	"github.com/fluent/fluent-logger-golang/fluent"
)

const (
	DEFAULT_APP_TAG_PREFIX = "service."

	EXT_TAG_PREFIX  = "EXT."
	EXT_TAG_PATTERN = "$$EXT$$:%s "
	EXT_TAG_REGEX   = "\\$\\$EXT\\$\\$:(\\w+)"

	ENV_PREFIX_KEY = "ENV_FLUENTD_PREFIX"
)

var fluentdLogger *fluent.Fluent
var extRegex = regexp.MustCompile(EXT_TAG_REGEX)

type FluentdWriter struct {
	defaultTag string // Default Fluentd Tag
}

func (fw *FluentdWriter) ReceiveMessage(message string, level seelog.LogLevel, context seelog.LogContextInterface) error {
	fluetdTag := fw.defaultTag
	if isExtMsg := extRegex.MatchString(message); isExtMsg == true {
		fluetdTag = EXT_TAG_PREFIX + getExtMessageTag(message)
	}

	var data = map[string]string{
		"app":     fluetdTag,
		"time":    context.CallTime().String(),
		"level":   level.String(),
		"message": message,
	}
	err := fluentdLogger.Post(fluetdTag, data)
	if err != nil {
		fmt.Println("ReceiveMessage error", err)
	}

	return nil
}

func (fw *FluentdWriter) AfterParse(initArgs seelog.CustomReceiverInitArgs) error {
	if len(os.Args) > 0 {
		envPrefix := beego.AppConfig.String(ENV_PREFIX_KEY)
		fw.defaultTag = DEFAULT_APP_TAG_PREFIX + envPrefix + filepath.Base(os.Args[0])
	} else {
		fw.defaultTag = "NoTag"
	}
	fmt.Printf("fluent tag is [%s]\r\n", fw.defaultTag)

	var err error
	fluentHost, fluentPort := getFluentEndpoint()
	fmt.Printf("fluentd address: %s:%d\r\n", fluentHost, fluentPort)
	fluentdLogger, err = fluent.New(fluent.Config{FluentPort: fluentPort, FluentHost: fluentHost, MaxRetry: -1})
	if err != nil {
		seelog.Errorf("connect to fluent failed: %s", err.Error())
	}
	return nil
}

func (fw *FluentdWriter) Flush() {

}

func (fw *FluentdWriter) Close() error {
	if fluentdLogger != nil {
		return fluentdLogger.Close()
	}
	return nil
}

// Should only be call if it's Extension Log Message
func getExtMessageTag(message string) string {
	foundMatch := extRegex.FindStringSubmatch(message)
	fmt.Printf("%q\r\n", foundMatch)
	//it returns the match string such as [$$EXT$$:Device Device]
	if len(foundMatch) == 2 {
		return foundMatch[1]
	}
	fmt.Printf("Invalid Extension Log Message as %s\r\n", message)
	return ""
}

func getFluentEndpoint() (string, int) {
	endPoint := beego.AppConfig.String("fluentdUrl")
	if endPoint != "" {
		pair := strings.Split(endPoint, ":")
		if len(pair) == 2 {
			port, err := strconv.Atoi(pair[1])
			if err == nil {
				return pair[0], port
			}
		}
	}
	return "127.0.0.1", 24230
}
