package log

import (
	"bytes"
	"errors"
)

type LogBuilder struct {
	log *Log
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{
		log: &Log{
			level:   InfoLevel,
			hooks:   NewHooks(),
			entries: newPool[Entry](),
			buffers: newPool[bytes.Buffer](),
		},
	}
}

func (that *LogBuilder) WithLevel(level Level) *LogBuilder {
	that.log.level = level
	return that
}

func (that *LogBuilder) WithExporter(exporter Exporter) *LogBuilder {
	that.log.exporter = exporter
	return that
}

func (that *LogBuilder) WithHook(hook Hook) *LogBuilder {
	that.log.AddHook(Levels.Keys(), hook)
	return that
}

func (that *LogBuilder) WithHookForLevels(hook Hook, levels []Level) *LogBuilder {
	that.log.AddHook(levels, hook)
	return that
}

func (that *LogBuilder) Build() (*Log, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}
	return that.log, nil
}

func (that *LogBuilder) checkRequiredFields() error {
	if that.log.exporter == nil {
		return ErrRequiredFieldExporter
	}

	return nil
}

var (
	ErrRequiredFieldExporter = errors.New("exporter is required")
)

func NewDummyLogger() *Log {
	l, err := NewLogBuilder().
		WithExporter(new(dummyExporter)).
		Build()
	if err != nil {
		panic(err)
	}
	return l
}
