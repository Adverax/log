package log

import "time"

const (
	DefaultTimestampFormat = time.RFC3339

	FieldKeyMsg         = "msg"
	FieldKeyLevel       = "level"
	FieldKeyTime        = "time"
	FieldKeyLoggerError = "logger_error"
	FieldKeyEntity      = "entity"
	FieldKeyAction      = "action"
	FieldKeyMethod      = "method"
	FieldKeySubject     = "subject"
	FieldKeyData        = "data"
	FieldKeyDuration    = "duration"
)

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

type fieldKey string

type FieldMap map[fieldKey]string

func (that FieldMap) Resolve(key fieldKey) string {
	if k, ok := that[key]; ok {
		return k
	}

	return string(key)
}

func (that FieldMap) PrefixFieldClashes(data Fields) {
	timeKey := that.Resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := that.Resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := that.Resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := that.Resolve(FieldKeyLoggerError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}
}
