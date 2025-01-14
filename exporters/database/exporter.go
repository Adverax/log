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
	data := that.makeData(entry)
	fields := that.extractFields(data)
	args := that.makeQueryArgs(fields)
	_, _ = that.db.Exec(that.query, args...)
}

func (that *Exporter) makeData(entry *log.Entry) log.Fields {
	data := entry.Data.Expand()

	that.fieldMap.EncodePrefixFieldClashes(data)

	timestampFormat := that.timestampFormat

	data[that.fieldMap.Resolve(log.FieldKeyTime)] = entry.Time.Format(timestampFormat)
	data[that.fieldMap.Resolve(log.FieldKeyMsg)] = entry.Message
	data[that.fieldMap.Resolve(log.FieldKeyLevel)] = entry.Level.String()

	return data
}

func (that *Exporter) extractFields(data log.Fields) log.Fields {
	fields := make(log.Fields, 4)
	for k, v := range data {
		if kk, ok := that.fieldMap[log.FieldKey(k)]; ok {
			fields[kk] = v
			delete(data, k)
		}
	}

	if that.dataKey != "" {
		raw, err := json.Marshal(data)
		if err == nil {
			fields[that.dataKey] = string(raw)
		}
	}

	return fields
}

func (that *Exporter) makeQueryArgs(fields log.Fields) []interface{} {
	args := make([]interface{}, 0, 32)
	for _, field := range that.fieldList {
		v := fields.Fetch(field)
		args = append(args, v)
	}
	return args
}
