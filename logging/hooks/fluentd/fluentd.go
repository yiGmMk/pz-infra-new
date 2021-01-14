package fluentd

import (
	"fmt"

	"github.com/yiGmMk/pz-infra-new/logging/hooks"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

var (
	FluentdTagName = "td"

	FluentdTagField     = "tag"
	FluentdTimeField    = "time"
	FluentdLevelField   = "level"
	FluentdMessageField = "message"
	fluentdFields       = []string{FluentdTagField, FluentdTimeField, FluentdLevelField, FluentdMessageField}
)

type fluentdHook struct {
	levels          []logrus.Level
	logger          *fluent.Fluent
	host            string
	port            int
	tag             string
	needReplaceDots bool
}

func New(host string, port int, tag string) (*fluentdHook, error) {
	h := &fluentdHook{
		levels: logrus.AllLevels,
		host:   host,
		port:   port,
		tag:    tag,
	}
	return h, h.connect()
}

func (h *fluentdHook) connect() error {
	logger, err := fluent.New(fluent.Config{
		FluentHost: h.host,
		FluentPort: h.port,
	})
	if err != nil {
		return err
	}
	h.logger = logger
	return nil
}

func (h *fluentdHook) Close() error {
	if h.logger != nil {
		return h.logger.Close()
	}
	return nil
}

func (h *fluentdHook) Levels() []logrus.Level {
	return h.levels
}

func (h *fluentdHook) SetLevels(levels []logrus.Level) *fluentdHook {
	h.levels = levels
	return h
}

// ReplaceDots replaces all `.` in field names and tags with `_` (an underscore).
// Because of field names cannot contain the `.` in Elasticsearch 2.0.
// This means you may need to change some conditionals in your Logstash configuration.
// See https://discuss.elastic.co/t/field-name-cannot-contain/33251.
func (h *fluentdHook) ReplaceDots() *fluentdHook {
	h.needReplaceDots = true
	return h
}

func (h *fluentdHook) Fire(entry *logrus.Entry) error {
	return h.logger.PostWithTime(h.tag, entry.Time, h.constructData(entry))
}

func (h *fluentdHook) constructData(entry *logrus.Entry) interface{} {
	solveFieldConflicts(entry.Data, fluentdFields, "fileds.")

	data := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		if !isReserved(k) {
			data[k] = v
		}
	}

	data[FluentdTagField] = h.tag
	data[FluentdTimeField] = entry.Time
	data[FluentdMessageField] = entry.Message
	data[FluentdLevelField] = hooks.LevelString(entry.Level)

	return ConvertToValue(data, FluentdTagName, h.needReplaceDots)
}

func isReserved(key string) bool {
	return key == "time" || key == "msg" || key == "level"
}

func solveFieldConflicts(data logrus.Fields, conflictingFields []string, prefix string) {
	var fixed string
	for _, f := range conflictingFields {
		if d, ok := data[f]; ok {
			fixed = fmt.Sprint(prefix, f)
			data[fixed] = d
			delete(data, f)
		}
	}
}
