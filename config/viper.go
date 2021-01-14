package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ViperProvider is a configuration provider.
type ViperProvider struct {
}

// New returns a Configer implemented by [Viper](https://github.com/spf13/viper).
func (vp *ViperProvider) New(option interface{}) (Configer, error) {
	opt, err := setOption(option)
	if err != nil {
		return nil, err
	}
	return newViperConfiger(opt), nil
}

type viperConfiger struct {
	Viper *viper.Viper
}

// ViperOption is used to set options for Viper.
type ViperOption struct {
	ConfigType string
	ConfigFile string
	IsWatch    bool
}

const (
	defaultConfigType = "yaml"
	defaultConfigFile = "config.yaml"
)

func setOption(option interface{}) (*ViperOption, error) {
	if option == nil {
		return nil, fmt.Errorf("viper: option is nil")
	}
	opt, ok := option.(*ViperOption)
	if !ok {
		return nil, fmt.Errorf("viper: the type of option should be (%s)", reflect.TypeOf(&ViperOption{}))
	}
	if opt.ConfigType == "" {
		opt.ConfigType = defaultConfigType
	}
	if opt.ConfigFile == "" {
		opt.ConfigFile = defaultConfigFile
	}
	return opt, nil
}

func newViperConfiger(option *ViperOption) *viperConfiger {
	vc := &viperConfiger{
		Viper: viper.New(),
	}
	vc.Viper.SetConfigType(option.ConfigType)
	vc.Viper.SetConfigFile(option.ConfigFile)
	fmt.Printf("viper: config file [%s] \n", vc.Viper.ConfigFileUsed())
	err := vc.Viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("viper: failed to read config file [%s] \n", err))
	}
	if option.IsWatch {
		vc.Viper.WatchConfig()
		fmt.Printf("viper: config file [%s] is watching \n", vc.Viper.ConfigFileUsed())

		vc.Viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("viper: config file [%s] changed \n", e.Name)
		})
	}
	return vc
}

func (vc *viperConfiger) Get(key string) interface{} {
	return vc.Viper.Get(key)
}

func (vc *viperConfiger) String(key string) string {
	return vc.Viper.GetString(key)
}

func (vc *viperConfiger) Bool(key string) bool {
	return vc.Viper.GetBool(key)
}

func (vc *viperConfiger) Int(key string) int {
	return vc.Viper.GetInt(key)
}

func (vc *viperConfiger) Float64(key string) float64 {
	return vc.Viper.GetFloat64(key)
}

func (vc *viperConfiger) Time(key string) time.Time {
	return vc.Viper.GetTime(key)
}

func (vc *viperConfiger) Map(key string) map[string]interface{} {
	return vc.Viper.GetStringMap(key)
}

func (vc *viperConfiger) Strings(key string) []string {
	return vc.Viper.GetStringSlice(key)
}

func (vc *viperConfiger) StringMap(key string) map[string]string {
	return vc.Viper.GetStringMapString(key)
}

func (vc *viperConfiger) Set(key string, value interface{}) {
	vc.Viper.Set(key, value)
}

func (vc *viperConfiger) IsSet(key string) bool {
	return vc.Viper.IsSet(key)
}
