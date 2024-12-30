package log

import "github.com/Adverax/core"

type FormatterJsonBuilder struct {
	*core.Builder
	formatter *JSONFormatter
}

func NewFormatterJsonBuilder() *FormatterJsonBuilder {
	return &FormatterJsonBuilder{
		Builder: core.NewBuilder("json_formatter"),
		formatter: &JSONFormatter{
			timestampFormat:   defaultTimestampFormat,
			disableTimestamp:  false,
			disableHTMLEscape: false,
			dataKey:           "",
			fieldMap:          nil,
			prettyPrint:       false,
		},
	}
}

func (that *FormatterJsonBuilder) WithDataKey(key string) *FormatterJsonBuilder {
	that.formatter.dataKey = key
	return that
}

func (that *FormatterJsonBuilder) WithFieldMap(fieldMap FieldMap) *FormatterJsonBuilder {
	that.formatter.fieldMap = fieldMap
	return that
}

func (that *FormatterJsonBuilder) WithPrettyPrint(prettyPrint bool) *FormatterJsonBuilder {
	that.formatter.prettyPrint = prettyPrint
	return that
}

func (that *FormatterJsonBuilder) WithTimestampFormat(timestampFormat string) *FormatterJsonBuilder {
	that.formatter.timestampFormat = timestampFormat
	return that
}

func (that *FormatterJsonBuilder) WithDisableTimestamp(disableTimestamp bool) *FormatterJsonBuilder {
	that.formatter.disableTimestamp = disableTimestamp
	return that
}

func (that *FormatterJsonBuilder) WithDisableHTMLEscape(disableHTMLEss bool) *FormatterJsonBuilder {
	that.formatter.disableHTMLEscape = disableHTMLEss
	return that
}

func (that *FormatterJsonBuilder) Build() (*JSONFormatter, error) {
	return that.formatter, nil
}
