package database

import (
	"database/sql"
	"errors"
	"github.com/adverax/core"
	"github.com/adverax/log"
	"strings"
)

type Builder struct {
	*core.Builder
	renderer *Renderer
}

func NewBuilder() *Builder {
	return &Builder{
		Builder: core.NewBuilder("database-renderer"),
		renderer: &Renderer{
			table:           "log",
			dataKey:         log.FieldKeyData,
			timestampFormat: log.DefaultTimestampFormat,
			fieldMap:        make(log.FieldMap),
		},
	}
}

func (that *Builder) WithDatabase(db *sql.DB) *Builder {
	that.renderer.db = db
	return that
}

func (that *Builder) WithTable(table string) *Builder {
	that.renderer.table = table
	return that
}

func (that *Builder) WithFieldMap(fieldMap log.FieldMap) *Builder {
	that.renderer.fieldMap = fieldMap
	return that
}

func (that *Builder) WithDataKey(dataKey string) *Builder {
	that.renderer.dataKey = dataKey
	return that
}

func (that *Builder) WithTimestampFormat(timestampFormat string) *Builder {
	that.renderer.timestampFormat = timestampFormat
	return that
}

func (that *Builder) Build() (*Renderer, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	that.renderer.fieldList = that.makeFieldList()
	that.renderer.query = that.makeQuery()
	return that.renderer, nil
}

func (that *Builder) checkRequiredFields() error {
	that.Builder.RequiredField(that.renderer.db, ErrRequiredFieldDatabase)
	that.RequiredField(that.renderer.table, ErrRequiredFieldTable)
	that.RequiredField(that.renderer.fieldMap, ErrRequiredFieldFieldMap)
	that.RequiredField(that.renderer.dataKey, ErrRequiredFieldDataKey)
	that.RequiredField(that.renderer.timestampFormat, ErrRequiredFieldTimestampFormat)

	return that.ResError()
}

func (that *Builder) makeFieldList() []string {
	fields := make([]string, 0, len(that.renderer.fieldMap))
	for _, v := range that.renderer.fieldMap {
		fields = append(fields, v)
	}
	return fields
}

func (that *Builder) makeQuery() string {
	return "INSERT INTO " + that.renderer.table + " (" + strings.Join(that.renderer.fieldList, ", ") + ") VALUES (" + strings.Repeat("?, ", len(that.renderer.fieldList)-1) + "?)"
}

var (
	ErrRequiredFieldDatabase        = errors.New("database is required")
	ErrRequiredFieldTable           = errors.New("table is required")
	ErrRequiredFieldFieldMap        = errors.New("field map is required")
	ErrRequiredFieldDataKey         = errors.New("data key is required")
	ErrRequiredFieldTimestampFormat = errors.New("timestamp format is required")
)
