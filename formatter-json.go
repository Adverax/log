package log

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
	timestampFormat   string
	disableTimestamp  bool
	disableHTMLEscape bool
	dataKey           string
	fieldMap          FieldMap
	prettyPrint       bool
}

// Format renders a single log entry
func (that *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := entry.Data.Expand()

	if that.dataKey != "" {
		newData := make(Fields, 4)
		newData[that.dataKey] = data
		data = newData
	}

	prefixFieldClashes(data, that.fieldMap)

	timestampFormat := that.timestampFormat

	if entry.err != "" {
		data[that.fieldMap.resolve(FieldKeyLoggerError)] = entry.err
	}
	if !that.disableTimestamp {
		data[that.fieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	}
	data[that.fieldMap.resolve(FieldKeyMsg)] = entry.Message
	data[that.fieldMap.resolve(FieldKeyLevel)] = entry.Level.String()

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!that.disableHTMLEscape)
	if that.prettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}
