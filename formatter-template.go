package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"text/template"
)

type Purifier interface {
	Purify(original, derivative string) string
}

// TemplateFormatter formats logs into text
type TemplateFormatter struct {
	purifier               Purifier
	disableTimestamp       bool
	timestampFormat        string
	disableSorting         bool
	sortingFunc            func([]string)
	disableLevelTruncation bool
	padLevelText           bool
	fieldMap               FieldMap
	template               *template.Template
	systemFields           map[string]struct{}
}

// Format renders a single log entry
func (that *TemplateFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	that.fieldMap.EncodePrefixFieldClashes(data)
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	fixedKeys := make([]string, 0, 4+len(data))
	if !that.disableTimestamp {
		fixedKeys = append(fixedKeys, that.fieldMap.Resolve(FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, that.fieldMap.Resolve(FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, that.fieldMap.Resolve(FieldKeyMsg))
	}
	if entry.err != "" {
		fixedKeys = append(fixedKeys, that.fieldMap.Resolve(FieldKeyLoggerError))
	}

	if !that.disableSorting {
		if that.sortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			fixedKeys = append(fixedKeys, keys...)
			that.sortingFunc(fixedKeys)

		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := that.timestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	systemFields := that.systemFields
	if systemFields == nil {
		systemFields = defaultSystemFields
	}

	params := make(map[string]interface{})
	rest := make(map[string]interface{})
	var entity, action, method string
	var subject, body string

	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case key == that.fieldMap.Resolve(FieldKeyTime):
			value = entry.Time.Format(timestampFormat)
		case key == that.fieldMap.Resolve(FieldKeyLevel):
			value = entry.Level.String()
		case key == that.fieldMap.Resolve(FieldKeyMsg):
			value = that.purify(entry.Message)
		case key == that.fieldMap.Resolve(FieldKeyLoggerError):
			value = entry.err
		case key == that.fieldMap.Resolve(FieldKeyTraceID):
			value, _ = data[key]
			continue
		case key == that.fieldMap.Resolve(FieldKeyEntity):
			value, _ = data[key]
			entity = that.value2string(value)
			continue
		case key == that.fieldMap.Resolve(FieldKeyAction):
			value, _ = data[key]
			action = that.value2string(value)
			continue
		case key == that.fieldMap.Resolve(FieldKeyMethod):
			value, _ = data[key]
			method = that.value2string(value)
			continue
		case key == that.fieldMap.Resolve(FieldKeySubject):
			value, _ = data[key]
			subject = that.value2string(value)
			continue
		case key == that.fieldMap.Resolve(FieldKeyData):
			value, _ = data[key]
			body = that.value2string(value)
			continue
		default:
			value = data[key]
		}

		val := that.value2string(value)
		if _, ok := systemFields[key]; ok {
			params[key] = val
		} else {
			rest[key] = val
		}
	}

	params["entity"] = that.formatEntity(entity, action)
	params["event"] = that.formatEvent(method, subject, body)

	if _, ok := params[FieldKeyTraceID]; !ok {
		params[FieldKeyTraceID] = ""
	}

	if len(rest) > 0 {
		var details []byte
		details, _ = json.Marshal(rest)
		params["details"] = string(details)
	}

	tpl := that.template
	if tpl == nil {
		tpl = defaultTpl
	}
	_ = tpl.Execute(b, params)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (that *TemplateFormatter) value2string(value interface{}) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	return stringVal
}

func (that *TemplateFormatter) formatEntity(entity, action string) string {
	if entity == "" {
		return ""
	}

	var result bytes.Buffer
	result.WriteByte(' ')
	result.WriteString(entity)
	if action == "" {
		result.WriteByte(':')
	} else {
		result.WriteByte(' ')
		result.WriteString(action)
	}

	return result.String()
}

func (that *TemplateFormatter) formatEvent(method, subject, body string) string {
	if body == "" && subject == "" {
		return ""
	}

	var wantSpace bool
	var result bytes.Buffer
	result.WriteByte(' ')

	if method != "" {
		result.WriteString(method)
		wantSpace = true
	}

	if subject != "" {
		if wantSpace {
			result.WriteByte(' ')
		}
		result.WriteString(subject)
		wantSpace = true
	}

	if body != "" {
		if wantSpace {
			result.WriteByte(' ')
		}
		result.WriteString(that.purify(body))
	}

	return result.String()
}

func (that *TemplateFormatter) purify(s string) string {
	if that.purifier == nil {
		return s
	}

	return that.purifier.Purify(s, s)
}
