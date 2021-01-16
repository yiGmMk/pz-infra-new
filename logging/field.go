package logging

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

// A Field is used to add a key-value pair to a logger's context.
type Field struct {
	key   string
	value interface{}
}

func (f Field) Key() string {
	return f.key
}

func (f Field) Value() interface{} {
	return f.value
}

// With constructs a Field with the given key and value.
func With(key string, value interface{}) Field {
	return Field{
		key:   key,
		value: value,
	}
}

// With constructs a Field with error.
func WithError(err error) Field {
	return Field{
		key:   logrus.ErrorKey,
		value: err,
	}
}

// 以%+v格式输出error,方便输出错误的详细信息
func FormatError(err error, args ...interface{}) string {
	return fmt.Sprintln(fmt.Sprintf("%+v", err), args)
}

// Stacktrace constructs a Field that stores a stacktrace of the current goroutine
// under the key "stacktrace".
func Stacktrace() Field {
	stacktrace := takeStacktrace(false)
	return With("stacktrace", stacktrace)
}

func takeStacktrace(includeAllGoroutines bool) string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, includeAllGoroutines)
	for n >= len(buf) {
		// Buffer wasn't large enough, allocate a larger one. No need to copy
		// previous buffer's contents.
		size := 2 * n
		buf = make([]byte, size)
		n = runtime.Stack(buf, includeAllGoroutines)
	}
	return string(buf[:n])
}
