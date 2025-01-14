package rex

import (
	"encoding/json"
	"github.com/adverax/log"
	"strings"
	"time"
)

type Engine struct {
	frames           []Frame
	fieldMap         log.FieldMap
	disableTimestamp bool
	timestampFormat  string
}

func (that *Engine) Parse(data []byte, entry *log.Entry) (int, error) {
	fields, n, err := that.parse(data)
	if err != nil {
		return n, err
	}

	that.consume(fields, entry)
	that.fieldMap.DecodePrefixFieldClashes(entry.Data)

	return n, nil
}

func (that *Engine) parse(data []byte) (fields map[string]string, nn int, err error) {
	fields = make(map[string]string)
	for _, frame := range that.frames {
		n, err := frame.Parse(data, fields)
		if err != nil {
			return nil, n, err
		}
		nn += n
		data = data[n:]
	}
	return fields, nn, nil
}

func (that *Engine) consume(fields map[string]string, entry *log.Entry) {
	if !that.disableTimestamp {
		key := that.fieldMap.Resolve(log.FieldKeyTime)
		if t, ok := fields[key]; ok {
			v, err := time.Parse(that.timestampFormat, t)
			if err == nil {
				entry.Time = v
			}
		}
		delete(fields, key)
	}

	key := that.fieldMap.Resolve(log.FieldKeyMsg)
	if v, ok := fields[key]; ok {
		entry.Message = v
		delete(fields, key)
	}

	key = that.fieldMap.Resolve(log.FieldKeyLevel)
	if v, ok := fields[key]; ok {
		level, _ := log.Levels.Encode(strings.ToLower(v))
		entry.Level = level
		delete(fields, key)
	}

	key = that.fieldMap.Resolve(log.FieldKeyData)
	if v, ok := fields[key]; ok {
		_ = json.Unmarshal([]byte(v), &entry.Data)
		delete(fields, key)
	}

	for k, v := range fields {
		entry.Data[k] = v
	}
}
