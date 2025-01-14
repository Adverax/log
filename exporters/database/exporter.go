package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/adverax/log"
)

type Translator func(interface{}) interface{}

type Exporter struct {
	db              *sql.DB
	table           string
	fieldMap        log.FieldMap
	query           string
	dataKey         string
	timestampFormat string
	fieldList       []string
}

func (that *Exporter) Export(ctx context.Context, entry *log.Entry) {
	data := entry.Data.Expand()

	that.fieldMap.EncodePrefixFieldClashes(data)

	timestampFormat := that.timestampFormat

	data[that.fieldMap.Resolve(log.FieldKeyTime)] = entry.Time.Format(timestampFormat)
	data[that.fieldMap.Resolve(log.FieldKeyMsg)] = entry.Message
	data[that.fieldMap.Resolve(log.FieldKeyLevel)] = entry.Level.String()

	fields := make(log.Fields, 4)
	for k, v := range data {
		if kk, ok := that.fieldMap[log.FieldKey(k)]; ok {
			fields[kk] = v
			delete(data, k)
		}
	}

	if that.dataKey != "" {
		raw, err := json.Marshal(entry.Data)
		if err == nil {
			fields[that.dataKey] = string(raw)
		}
	}

	args := make([]interface{}, 0, 32)
	for _, field := range that.fieldList {
		v := fields.Fetch(field)
		args = append(args, v)
	}

	_, _ = that.db.Exec(that.query, args...)
}
