package log

import (
	"errors"
	"strings"
	"text/template"
)

type FormatterTemplateBuilder struct {
	formatter *TemplateFormatter
}

func NewFormatterTemplateBuilder() *FormatterTemplateBuilder {
	return &FormatterTemplateBuilder{
		formatter: &TemplateFormatter{
			timestampFormat:  DefaultTimestampFormat,
			disableTimestamp: false,
			template:         defaultTpl,
			systemFields:     defaultSystemFields,
			fieldMap:         FieldMap{},
		},
	}
}

func (that *FormatterTemplateBuilder) WithTemplate(tpl *template.Template) *FormatterTemplateBuilder {
	that.formatter.template = tpl
	return that
}

func (that *FormatterTemplateBuilder) WithDisableTimestamp(disableTimestamp bool) *FormatterTemplateBuilder {
	that.formatter.disableTimestamp = disableTimestamp
	return that
}

func (that *FormatterTemplateBuilder) WithTimestampFormat(timestampFormat string) *FormatterTemplateBuilder {
	that.formatter.timestampFormat = timestampFormat
	return that
}

func (that *FormatterTemplateBuilder) WithFieldMap(fieldMap FieldMap) *FormatterTemplateBuilder {
	that.formatter.fieldMap = fieldMap
	return that
}

func (that *FormatterTemplateBuilder) WithDisableSorting(disableSorting bool) *FormatterTemplateBuilder {
	that.formatter.disableSorting = disableSorting
	return that
}

func (that *FormatterTemplateBuilder) WithSortingFunc(sortingFunc func([]string)) *FormatterTemplateBuilder {
	that.formatter.sortingFunc = sortingFunc
	return that
}

func (that *FormatterTemplateBuilder) WithDisableLevelTruncation(disableLevelTruncation bool) *FormatterTemplateBuilder {
	that.formatter.disableLevelTruncation = disableLevelTruncation
	return that
}

func (that *FormatterTemplateBuilder) WithPadLevelText(padLevelText bool) *FormatterTemplateBuilder {
	that.formatter.padLevelText = padLevelText
	return that
}

func (that *FormatterTemplateBuilder) Build() (*TemplateFormatter, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.formatter, nil
}

func (that *FormatterTemplateBuilder) checkRequiredFields() error {
	if that.formatter.template == nil {
		return ErrTemplateRequired
	}

	return nil
}

var (
	ErrTemplateRequired = errors.New("template is required")
)

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}

var defaultTemplate = `{{.time}} {{.level | ToUpper}}{{if .trace_id}} #{{.trace_id}}{{end}}:{{.entity}} {{.msg}}{{.event}}{{if .details}} DETAILS {{.details}}{{end}}`

var defaultTpl = template.Must(template.New("log").Funcs(funcMap).Parse(defaultTemplate))

var defaultSystemFields = map[string]struct{}{
	FieldKeyTime:    {},
	FieldKeyLevel:   {},
	FieldKeyMsg:     {},
	FieldKeyTraceID: {},
	FieldKeyEntity:  {},
	FieldKeyAction:  {},
	FieldKeyMethod:  {},
	FieldKeySubject: {},
	FieldKeyData:    {},
}
