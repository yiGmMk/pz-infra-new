package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

const Logrus = "logrus"
const skipCallerDepth = 3

// LogrusProvider is a logger provider.
type LogrusProvider struct {
}

// New returns a Logger implemented by [Logrus](https://github.com/sirupsen/logrus).
func (lp *LogrusProvider) New(option interface{}) (Logger, error) {
	opt, err := setOption(option)
	if err != nil {
		return nil, err
	}
	return newLogrusLogger(opt), nil
}

type logrusLogger struct {
	Component  string
	Logrus     *logrus.Logger
	workingDir string
}

// LogrusOption is used to set options for Logrus.
type LogrusOption struct {
	Out       io.Writer
	Hooks     []logrus.Hook
	Formatter logrus.Formatter
	Level     logrus.Level
	Component string
}

var (
	defaultLogrusOption = &LogrusOption{
		Out:       os.Stderr,
		Hooks:     []logrus.Hook{},
		Formatter: &logrus.TextFormatter{},
		Level:     logrus.DebugLevel,
	}
)

func (o *LogrusOption) clone() *LogrusOption {
	copy := *o
	return &copy
}

func (o *LogrusOption) levelHooks() logrus.LevelHooks {
	lh := make(logrus.LevelHooks)
	for _, hook := range o.Hooks {
		fmt.Println(hook.Levels())
		for _, level := range hook.Levels() {
			lh[level] = append(lh[level], hook)
		}
	}
	return lh
}

func setOption(option interface{}) (*LogrusOption, error) {
	if option == nil {
		return nil, fmt.Errorf("logrus: option is nil")
	}
	opt, ok := option.(*LogrusOption)
	if !ok {
		return nil, fmt.Errorf("logrus: the type of option should be (%s)", reflect.TypeOf(&LogrusOption{}))
	}
	newOpt := defaultLogrusOption.clone()
	if opt.Out != nil {
		newOpt.Out = opt.Out
	}
	if len(opt.Hooks) > 0 {
		newOpt.Hooks = opt.Hooks
	}
	if opt.Formatter != nil {
		newOpt.Formatter = opt.Formatter
	}
	newOpt.Level = opt.Level
	newOpt.Component = opt.Component
	return newOpt, nil
}

func newLogrusLogger(option *LogrusOption) *logrusLogger {
	return &logrusLogger{
		Component: option.Component,
		Logrus: &logrus.Logger{
			Out:       option.Out,
			Formatter: option.Formatter,
			Hooks:     option.levelHooks(),
			Level:     option.Level,
		},
		workingDir: getWorkingDir(),
	}
}

func getWorkingDir() string {
	workingDir := "/"
	wd, err := os.Getwd()
	if err == nil {
		workingDir = filepath.ToSlash(wd) + "/"
	}
	return workingDir
}

func (log *logrusLogger) addFields(fields []Field, needStackTrace bool) logrus.Fields {
	fs := logrus.Fields{}
	for _, f := range fields {
		fs[f.Key()] = f.Value()
	}
	if log.Component != "" {
		fs["component"] = log.Component
	}
	if needStackTrace {
		stackTrace := Stacktrace()
		fs[stackTrace.Key()] = stackTrace.Value()
	}
	fp, _, fn, ln, _ := log.extractCallerInfo(skipCallerDepth)
	_, fileName := filepath.Split(fp)
	fs["fileFunc"] = fmt.Sprintf("%s:%d:%s", fileName, ln, fn)
	return fs
}

func (log *logrusLogger) extractCallerInfo(skip int) (fullPath string, shortPath string, funcName string, line int, err error) {
	pc, fp, ln, ok := runtime.Caller(skip)
	if !ok {
		err = fmt.Errorf("error during runtime.Caller")
		return
	}
	line = ln
	fullPath = fp
	if strings.HasPrefix(fp, log.workingDir) {
		shortPath = fp[len(log.workingDir):]
	} else {
		shortPath = fp
	}
	funcNameFull := runtime.FuncForPC(pc).Name()
	if strings.HasPrefix(funcNameFull, log.workingDir) {
		funcNameFull = funcNameFull[len(log.workingDir):]
	}
	funcSplit := strings.Split(funcNameFull, ".")
	if len(funcSplit) > 0 {
		funcName = funcSplit[len(funcSplit)-1]
	}
	return
}

func (log *logrusLogger) Debug(message string, fields ...Field) {
	log.Logrus.WithFields(log.addFields(fields, false)).Debugln(message)
}

func (log *logrusLogger) Info(message string, fields ...Field) {
	log.Logrus.WithFields(log.addFields(fields, false)).Infoln(message)
}

func (log *logrusLogger) Warn(message string, fields ...Field) {
	log.Logrus.WithFields(log.addFields(fields, false)).Warnln(message)
}

func (log *logrusLogger) Error(message string, fields ...Field) error {
	log.Logrus.WithFields(log.addFields(fields, true)).Errorln(message)
	return errors.New(message)
}

func (log *logrusLogger) Fatal(message string, fields ...Field) {
	log.Logrus.WithFields(log.addFields(fields, true)).Fatalln(message)

}

func (log *logrusLogger) Panic(message string, fields ...Field) {
	log.Logrus.WithFields(log.addFields(fields, true)).Panicln(message)
}

func (log *logrusLogger) GetPrintLogger() PrintLogger {
	return printLogger{Logrus: log.Logrus}
}

type printLogger struct {
	Logrus *logrus.Logger
}

func (l printLogger) Print(v ...interface{}) {
	isError := false
	for _, item := range v {
		_, ok := item.(error)
		if ok {
			isError = true
			break
		}
	}
	if isError {
		l.Logrus.Error(v...)
		l.Logrus.Error(Stacktrace().Value().(string))
	} else {
		l.Logrus.Debug(v...)
	}
}

func (l printLogger) Write(p []byte) (n int, err error) {
	l.Logrus.Info(string(p))
	return 0, nil
}
