package rolling

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/yiGmMk/pz-infra-new/logging/hooks"

	"github.com/sirupsen/logrus"
)

type LevelPaths map[logrus.Level]string
type levelFiles map[logrus.Level]*lumberjack.Logger
type pathFiles map[string]*lumberjack.Logger

var files = make(pathFiles)

// Hook to handle writing to rolling log files.
type rollingHook struct {
	levels []logrus.Level
	path   string
	file   *lumberjack.Logger
	paths  LevelPaths
	files  levelFiles
}

func New(path string) *rollingHook {
	return NewWithLevelPaths(path, nil)
}

func NewWithLevelPaths(path string, levelPaths LevelPaths) *rollingHook {
	hook := &rollingHook{
		levels: logrus.AllLevels,
		file:   newRollingFile(path),
	}
	if len(levelPaths) > 0 {
		hook.paths = levelPaths
		hook.files = make(levelFiles)
		for l, p := range levelPaths {
			hook.files[l] = newRollingFile(p)
		}
	}
	return hook
}

func newRollingFile(path string) *lumberjack.Logger {
	if path == "" {
		panic("rolling: log file path is empty")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("rolling: can't get absolute path for log file: %v", err))
	}
	file, ok := files[abs]
	if ok {
		return file
	}
	newFile := &lumberjack.Logger{
		Filename:   abs,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
	files[abs] = newFile
	return newFile
}

func (hook *rollingHook) Fire(entry *logrus.Entry) error {
	serialized := serialize(entry, time.RFC3339)
	_, err := hook.file.Write(serialized)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}
	defer hook.file.Close()

	f, ok := hook.files[entry.Level]
	if ok {
		_, err = f.Write(serialized)
		if err != nil {
			return fmt.Errorf("failed to write message: %v", err)
		}
		defer f.Close()
	}
	return nil
}

func (hook *rollingHook) Levels() []logrus.Level {
	return hook.levels
}

func serialize(entry *logrus.Entry, timestampFormat string) []byte {
	b := &bytes.Buffer{}
	levelString := strings.ToUpper(hooks.LevelString(entry.Level))
	fmt.Fprintf(b, "%-5s[%s] %-44s ", levelString, entry.Time.Format(timestampFormat), entry.Message)
	for k, v := range entry.Data {
		fmt.Fprintf(b, " %s=%+v", k, v)
	}
	b.WriteByte('\n')
	return b.Bytes()
}
