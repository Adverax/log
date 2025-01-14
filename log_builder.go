package log

import (
	"bytes"
	"github.com/adverax/core"
	"os"
)

type LogBuilder struct {
	*core.Builder
	log *Log
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{
		Builder: core.NewBuilder("log"),
		log: &Log{
			level:   InfoLevel,
			hooks:   NewHooks(),
			entries: core.NewPool[Entry](),
			buffers: core.NewPool[bytes.Buffer](),
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
	that.log.AddHook(Levels.Values(), hook)
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
	if err := that.updateDefaultFields(); err != nil {
		return nil, err
	}
	return that.log, nil
}

func (that *LogBuilder) checkRequiredFields() error {
	return that.ResError()
}

func (that *LogBuilder) updateDefaultFields() error {
	if that.log.exporter == nil {
		formatter, err := NewFormatterJsonBuilder().Build()
		if err != nil {
			return err
		}
		that.log.exporter = NewRenderer(formatter, os.Stdout)
	}

	return nil
}

func NewDummyLogger() *Log {
	return core.Must(
		NewLogBuilder().
			WithExporter(new(dummyExporter)).
			Build(),
	)
}
