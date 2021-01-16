package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Logger
type Logger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	// 无堆栈信息错误使用,已携带堆栈信息的error不应使用该函数,会输出双份的堆栈信息
	Error(message string, fields ...Field) error
	// Fatal logs a message, then calls os.Exit(1).
	Fatal(message string, fields ...Field)
	// Panic logs a message, then panics.
	Panic(message string, fields ...Field)
	// Get Print Method Logger
	GetPrintLogger() PrintLogger

	WithContext(ctx context.Context) *logrus.Entry

	// github.com/pkg/errors
	// errors.New/WithStack/Wrap/Wrapf等可以方便地给错误带上堆栈信息,方便定位错误
	// 已携带堆栈信息的错误,使用fmt.Sprintf("%+v",err)得到完整的错误信息输出
	// 未包含堆栈信息的错误请使用 Error()输出日志,会带上堆栈信息
	LogErrorHasStackInfo(err error, args ...interface{})
}

type PrintLogger interface {
	Print(v ...interface{})
	Write(p []byte) (n int, err error)
}

// Ensure to call InitLogger("componentName") first before using Log
var Log Logger
