package confUtil

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gyf841010/pz-infra-new/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

var (
	loadedConfigFileName, loadedAppName string
)

const (
	ENV_PREFIX      = "Z_"
	DEFAULT_SECTION = "default"
)

type InitBeegoConfigOptions struct {
	RenderedConfigFileName       string // file name to save rendered template config file, the file is saved at the same directory of template file
	AppName                      string // application name for config section override
	ConfigTemplateFileName       string // config template file name, default is app.conf.templ
	ConfigTemplateValuesFileName string // config template values file name, default is app.conf.values
	// file name of conf to load at final. default is the same as "RenderedConfigFileName" field,
	// use different name when necessary. In unit test init function, it use apps_test.conf
	ConfigFileNameToLoad string
}

func init() {
	// try to initialize beego as soon as possible
	// InitBeegoConfig(&InitBeegoConfigOptions{RenderedConfigFileName: "app.conf", AppName: filepath.Base(os.Args[0])})
}

func PanicInitBeegoConfigWithDefaultAppName(configFileName string) {
	if len(os.Args) < 1 {
		panic("Cannot get command line name")
	}
	PanicInitBeegoConfig(configFileName, filepath.Base(os.Args[0]))
}

func PanicInitBeegoConfig(configFileName, appName string) {
	log.Debug("appName is", appName)
	if err := InitBeegoConfig(&InitBeegoConfigOptions{RenderedConfigFileName: configFileName, AppName: appName}); err != nil {
		panic(err)
	}
}

func InitBeegoConfig(options *InitBeegoConfigOptions) error {
	if options == nil {
		return errors.New("options is nil")
	}
	if options.RenderedConfigFileName == "" || options.AppName == "" {
		return errors.New(fmt.Sprintf("Invalid options: %+v", *options))
	}
	if options.ConfigTemplateFileName == "" {
		options.ConfigTemplateFileName = "app.conf.templ"
	}
	if options.ConfigTemplateValuesFileName == "" {
		options.ConfigTemplateValuesFileName = "app.conf.values"
	}
	if options.ConfigFileNameToLoad == "" {
		options.ConfigFileNameToLoad = options.RenderedConfigFileName
	}

	if loadedConfigFileName == options.RenderedConfigFileName && loadedAppName == options.AppName {
		log.Debug("config already loaded")
		return nil
	}

	_, err := renderConfigTemplate(options)
	if err != nil {
		return err
	}
	renderedFilePath := FileHierarchyFind(options.RenderedConfigFileName)
	configContainer, err := config.NewConfig("ini", renderedFilePath)
	if err != nil {
		return log.Error(err.Error())
	}
	if err := overrideByEnvironmentVariables(configContainer); err != nil {
		return err
	}
	if err := configContainer.SaveConfigFile(renderedFilePath); err != nil {
		return log.Error(err.Error())
	}

	configFilePath := FileHierarchyFind(options.ConfigFileNameToLoad)
	if configFilePath == "" {
		return errors.New(fmt.Sprintf("not found %s in current or parent directories", options.ConfigFileNameToLoad))
	}

	envRunMode := os.Getenv("BEEGO_RUNMODE")
	if options.AppName != "" { // make beego use <appName> section as part of default section
		os.Setenv("BEEGO_RUNMODE", options.AppName)
		defer os.Setenv("BEEGO_RUNMODE", envRunMode)
	}
	beego.LoadAppConfig("ini", configFilePath)
	defer recoverRunMode(envRunMode)
	if options.AppName != "" {
		// merge <appName> section to default section
		beego.BConfig.RunMode = DEFAULT_SECTION
		section, err := beego.AppConfig.GetSection(options.AppName)
		if err == nil {
			for key, value := range section {
				if err := beego.AppConfig.Set(key, value); err != nil {
					log.Error("override default config failed", err.Error())
				}
			}
		} else {
			log.Warn("not found config section. ", err.Error())
		}
	}
	loadedConfigFileName = options.ConfigFileNameToLoad
	loadedAppName = options.AppName
	return nil
}

func OverrideBeegoConfig(configFilePath string) error {
	generateFilePath := configFilePath + ".generated"
	configContainer, err := config.NewConfig("ini", configFilePath)
	if err != nil {
		return log.Error(err.Error())
	}
	if err := overrideByEnvironmentVariables(configContainer); err != nil {
		return err
	}
	if err := configContainer.SaveConfigFile(generateFilePath); err != nil {
		return log.Error(err.Error())
	}
	if err := beego.LoadAppConfig("ini", generateFilePath); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func overrideByEnvironmentVariables(configContainer config.Configer) error {
	section, err := configContainer.GetSection(DEFAULT_SECTION)
	if err != nil {
		return log.Errorf("Failed to get default config section:%s", err.Error())
	}
	for _, item := range os.Environ() {
		pair := strings.SplitN(item, "=", 2)
		if strings.HasPrefix(strings.ToUpper(pair[0]), ENV_PREFIX) {
			key := pair[0][len(ENV_PREFIX):]
			log.Infof("Override config with %s=%s", key, pair[1])
			section[strings.ToLower(key)] = pair[1]
			if key == "runmode" {
				beego.BConfig.RunMode = pair[1]
			}
		}
	}
	return nil
}

func recoverRunMode(envRunMode string) {
	if envRunMode != "" {
		beego.BConfig.RunMode = envRunMode
	} else {
		if runmode := beego.AppConfig.String("RunMode"); runmode != "" {
			beego.BConfig.RunMode = runmode
		}
	}
}

func renderConfigTemplate(options *InitBeegoConfigOptions) (string, error) {
	templateFilePath := FileHierarchyFind(options.ConfigTemplateFileName)
	if templateFilePath == "" {
		return "", log.Errorf("not found %s in current or parent directories", options.ConfigTemplateFileName)
	}
	tempBytes, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		return "", log.Error(err)
	}
	templateValuePath := FileHierarchyFind(options.ConfigTemplateValuesFileName)
	if templateValuePath == "" {
		log.Infof("Not found %s in current or parent directories, skip config render", templateValuePath)
		return saveRenderedConfig(templateFilePath, string(tempBytes), options.RenderedConfigFileName)
	}
	valBytes, err := ioutil.ReadFile(templateValuePath)
	if err != nil {
		return "", log.Error(err)
	}
	if len(valBytes) <= 0 || len(tempBytes) <= 0 {
		return saveRenderedConfig(templateFilePath, string(tempBytes), options.RenderedConfigFileName)
	}
	parser := &config.IniConfig{}
	configContainer, err := parser.ParseData(valBytes)
	if err != nil {
		return "", log.Error(err)
	}
	valMap, err := configContainer.GetSection(DEFAULT_SECTION)
	if err != nil {
		return "", log.Error(err)
	}
	// change {{XXX}} to {{conf "XXX"}} to let tempalte engine call conf functiont to get value
	templateStr := strings.Replace(string(tempBytes), "{{", "{{conf \"", -1)
	templateStr = strings.Replace(templateStr, "}}", "\"}}", -1)
	funcMap := template.FuncMap{
		"conf": func(key string) (string, error) {
			if value, found := valMap[strings.ToLower(strings.TrimSpace(key))]; found {
				return value, nil
			}
			return "", errors.New(fmt.Sprintf("not found template value with key [%s] in file [%s]", key, templateValuePath))
		},
	}
	templ, err := template.New("").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", log.Error(err)
	}
	buffer := &bytes.Buffer{}
	if err := templ.Execute(buffer, nil); err != nil {
		return "", log.Error(err)
	}
	return saveRenderedConfig(templateFilePath, buffer.String(), options.RenderedConfigFileName)
}

func saveRenderedConfig(templateFilePath, configContent, fileName string) (string, error) {
	filePath := path.Join(path.Dir(templateFilePath), fileName)
	log.Infof("Save rendered config to %s", filePath)
	return filePath, ioutil.WriteFile(filePath, []byte(configContent), 0600)
}
