package log

import "time"

const (
	defaultTimestampFormat = time.RFC3339

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

func (that FieldMap) resolve(key fieldKey) string {
	if k, ok := that[key]; ok {
		return k
	}

	return string(key)
}

func prefixFieldClashes(data Fields, fieldMap FieldMap) {
	timeKey := fieldMap.resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := fieldMap.resolve(FieldKeyLoggerError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}
}
