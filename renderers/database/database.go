package database

import (
	"context"
	"database/sql"
	"github.com/adverax/log"
)

type Translator func(interface{}) interface{}

type Renderer struct {
	db              *sql.DB
	table           string // log -> db
	fieldMap        log.FieldMap
	query           string
	dataKey         string
	timestampFormat string
}

func (that *Renderer) Render(ctx context.Context, entry *log.Entry) {
	data := entry.Data.Expand()

	if that.dataKey != "" {
		newData := make(log.Fields, 4)
		newData[that.dataKey] = data
		data = newData
	}

	that.fieldMap.PrefixFieldClashes(data)

	timestampFormat := that.timestampFormat

	data[that.fieldMap.Resolve(log.FieldKeyTime)] = entry.Time.Format(timestampFormat)
	data[that.fieldMap.Resolve(log.FieldKeyMsg)] = entry.Message
	data[that.fieldMap.Resolve(log.FieldKeyLevel)] = entry.Level.String()

	args := make([]interface{}, 0, 32)
	for _, v := range that.fieldMap {
		vv := data.Fetch(v)
		args = append(args, vv)
	}

	_, _ = that.db.Exec(that.query, args...)
}
