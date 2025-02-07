package database

import (
	"database/sql"
	"errors"
	"github.com/adverax/log"
	"strings"
)

type Builder struct {
	exporter *Exporter
}

func NewBuilder() *Builder {
	return &Builder{
		exporter: &Exporter{
			table:           "log",
			dataKey:         log.FieldKeyData,
			timestampFormat: log.DefaultTimestampFormat,
			fieldMap:        make(log.FieldMap),
		},
	}
}

func (that *Builder) WithDatabase(db *sql.DB) *Builder {
	that.exporter.db = db
	return that
}

func (that *Builder) WithTable(table string) *Builder {
	that.exporter.table = table
	return that
}

func (that *Builder) WithFieldMap(fieldMap log.FieldMap) *Builder {
	that.exporter.fieldMap = fieldMap
	return that
}

func (that *Builder) WithDataKey(dataKey string) *Builder {
	that.exporter.dataKey = dataKey
	return that
}

func (that *Builder) WithTimestampFormat(timestampFormat string) *Builder {
	that.exporter.timestampFormat = timestampFormat
	return that
}

func (that *Builder) Build() (*Exporter, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	that.exporter.fieldList = that.makeFieldList()
	that.exporter.query = that.makeQuery()
	return that.exporter, nil
}

func (that *Builder) checkRequiredFields() error {
	if that.exporter.db == nil {
		return ErrRequiredFieldDatabase
	}
	if that.exporter.table == "" {
		return ErrRequiredFieldTable
	}
	if len(that.exporter.fieldMap) == 0 {
		return ErrRequiredFieldFieldMap
	}
	if that.exporter.timestampFormat == "" {
		return ErrRequiredFieldTimestampFormat
	}
	return nil
}

func (that *Builder) makeFieldList() []string {
	fields := make([]string, 0, len(that.exporter.fieldMap))
	for _, v := range that.exporter.fieldMap {
		fields = append(fields, v)
	}
	return fields
}

func (that *Builder) makeQuery() string {
	return "INSERT INTO " + that.exporter.table + " (" + strings.Join(that.exporter.fieldList, ", ") + ") VALUES (" + strings.Repeat("?, ", len(that.exporter.fieldList)-1) + "?)"
}

var (
	ErrRequiredFieldDatabase        = errors.New("database is required")
	ErrRequiredFieldTable           = errors.New("table is required")
	ErrRequiredFieldFieldMap        = errors.New("field map is required")
	ErrRequiredFieldTimestampFormat = errors.New("timestamp format is required")
)
