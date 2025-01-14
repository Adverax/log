package rex

import (
	"errors"
	"github.com/adverax/core"
	"github.com/adverax/log"
)

type Builder struct {
	*core.Builder
	engine *Engine
}

func NewBuilder() *Builder {
	return &Builder{
		Builder: core.NewBuilder("parser"),
		engine: &Engine{
			fieldMap:         log.FieldMap{},
			disableTimestamp: false,
			timestampFormat:  log.DefaultTimestampFormat,
		},
	}
}

func (that *Builder) WithFieldMap(fieldMap log.FieldMap) *Builder {
	that.engine.fieldMap = fieldMap
	return that
}

func (that *Builder) WithDisableTimestamp(disableTimestamp bool) *Builder {
	that.engine.disableTimestamp = disableTimestamp
	return that
}

func (that *Builder) WithTimestampFormat(timestampFormat string) *Builder {
	that.engine.timestampFormat = timestampFormat
	return that
}

func (that *Builder) WithFrame(frame ...Frame) *Builder {
	that.engine.frames = append(that.engine.frames, frame...)
	return that
}

func (that *Builder) Build() (*Engine, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	return that.engine, nil
}

func (that *Builder) checkRequiredFields() error {
	that.RequiredField(that.engine.frames, ErrFieldFramesRequired)
	that.RequiredField(that.engine.fieldMap, ErrFieldFieldMapRequired)

	return that.ResError()
}

var (
	ErrFieldFramesRequired   = errors.New("frames required")
	ErrFieldFieldMapRequired = errors.New("field map required")
)
