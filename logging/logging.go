package logging

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/gyf841010/pz-infra-new/logging/hooks/fluentd"
	"github.com/gyf841010/pz-infra-new/logging/hooks/rolling"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
)

var defaultCommonLogPath = "./log/common.log"
var defaultErrorLogPath = "./log/error.log"

// Provider is the interface that must be implemented by a logger provider.
type LogProvider interface {
	// New returns a new logger.
	// The option is a provider-specific value used to set option(s) for the logger.
	New(option interface{}) (Logger, error)
}

var (
	providersMu          sync.Mutex
	providers            = make(map[string]LogProvider)
	registerProviderOnce sync.Once
)

// Register makes a provider available by name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider LogProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("logging: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("logging: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Providers returns a sorted list of the names of the registered providers.
func Providers() []string {
	providersMu.Lock()
	defer providersMu.Unlock()
	var names []string
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetLogger returns a provided logger.
func GetLogger(name string, option interface{}) (Logger, error) {
	provider, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("config: unknown provider[%s] (forgot to Register?)", name)
	}
	return provider.New(option)
}

// Initialize Logger.
func InitLogger(componentName string) Logger {
	initProvider()
	var hooks []logrus.Hook
	rollingHook := rolling.NewWithLevelPaths(defaultCommonLogPath, rolling.LevelPaths{
		logrus.ErrorLevel: defaultErrorLogPath,
	})
	hooks = append(hooks, rollingHook)
	if fluentdHook := fluentd.BuildFluentdHook(componentName); fluentdHook != nil {
		hooks = append(hooks, fluentdHook)
	}
	logLevel := logrus.InfoLevel
	if beego.BConfig.RunMode == "dev" {
		logLevel = logrus.DebugLevel
	}
	logger, err := GetLogger(Logrus, &LogrusOption{
		Level:     logLevel,
		Hooks:     hooks,
		Component: componentName,
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
	})
	if err != nil {
		panic(fmt.Errorf("tests: failed to init logger: %v", err))
	}
	Log = logger
	return logger
}

func initProvider() {
	registerProviderOnce.Do(func() {
		Register(Logrus, &LogrusProvider{})
	})
	if beego.AppConfig.String("logger.defaultCommonLogPath") != "" {
		defaultCommonLogPath = beego.AppConfig.String("logger.defaultCommonLogPath")
	}
	if beego.AppConfig.String("logger.defaultCommonLogPath") != "" {
		defaultErrorLogPath = beego.AppConfig.String("logger.defaultErrorLogPath")
	}
}
