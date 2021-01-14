package log

import (
	"fmt"
	"runtime/debug"

	"github.com/cihub/seelog"
)

func Debug(v ...interface{}) {
	seelog.Debug(v...)
}

func Info(v ...interface{}) {
	seelog.Info(v...)
}

func Infof(format string, param ...interface{}) {
	seelog.Infof(format, param...)
}

func Warn(v ...interface{}) {
	seelog.Warn(v...)
}

func Warnf(format string, param ...interface{}) {
	seelog.Warnf(format, param...)
}

func Error(v ...interface{}) error {
	isError := false
	for _, item := range v {
		_, ok := item.(error)
		if ok {
			isError = true
			break
		}
	}
	err := seelog.Error(v...)
	if isError {
		seelog.Error(string(debug.Stack()))
	}
	return err
}

func Errorf(format string, params ...interface{}) error {
	isError := false
	for _, item := range params {
		_, ok := item.(error)
		if ok {
			isError = true
			break
		}
	}
	err := seelog.Errorf(format, params...)
	if isError {
		seelog.Error(string(debug.Stack()))
	}
	return err
}

func Trace(v ...interface{}) {
	seelog.Trace(v...)
}

func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

//Extension Log with pattern "$$EXT$$:Device .... "
func Ext(module string, originalMsg string) {
	prefixString := fmt.Sprintf(EXT_TAG_PATTERN, module)
	seelog.Info(prefixString + originalMsg)
}

//Extension Log with pattern "$$EXT$$:Device .... "
func Extf(module string, originalMsg string, param ...interface{}) {
	prefixString := fmt.Sprintf(EXT_TAG_PATTERN, module)
	seelog.Info(prefixString+originalMsg, param)
}

func Critical(v ...interface{}) error {
	isError := false
	for _, item := range v {
		_, ok := item.(error)
		if ok {
			isError = true
			break
		}
	}
	err := seelog.Critical(v...)
	if isError {
		seelog.Critical(string(debug.Stack()))
	}
	return err
}

func Criticalf(format string, params ...interface{}) error {
	isError := false
	for _, item := range params {
		_, ok := item.(error)
		if ok {
			isError = true
			break
		}
	}
	err := seelog.Criticalf(format, params...)
	if isError {
		seelog.Critical(string(debug.Stack()))
	}
	return err
}

func Flush() {
	seelog.Flush()
}

func InitWithLogFile(logfile string) {
	seelog.RegisterReceiver("fluentdReceiver", &FluentdWriter{})
	logger, err := seelog.LoggerFromConfigAsFile(logfile)
	if err == nil {
		logger.SetAdditionalStackDepth(1)
		seelog.ReplaceLogger(logger)
		seelog.Info("init log from config success.")
	} else {
		seelog.Error(err)
	}
}

type logger struct {
}

func (l *logger) Print(v ...interface{}) {
	seelog.Debug(v...)
}

func (l *logger) Write(p []byte) (n int, err error) {
	seelog.Info(string(p))
	return 0, nil
}

func CurrentLogger() *logger {
	return &logger{}
}
