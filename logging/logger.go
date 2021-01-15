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
	Error(message string, fields ...Field) error
	// Fatal logs a message, then calls os.Exit(1).
	Fatal(message string, fields ...Field)
	// Panic logs a message, then panics.
	Panic(message string, fields ...Field)
	// Get Print Method Logger
	GetPrintLogger() PrintLogger

	WithContext(ctx context.Context) *logrus.Entry
}

type PrintLogger interface {
	Print(v ...interface{})
	Write(p []byte) (n int, err error)
}

// Ensure to call InitLogger("componentName") first before using Log
var Log Logger
