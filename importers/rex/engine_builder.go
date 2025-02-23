package rexImporter

import (
	"errors"
	"github.com/adverax/log"
)

type Builder struct {
	engine *Engine
}

func NewBuilder() *Builder {
	return &Builder{
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
	if that.engine.frames == nil {
		return ErrFieldFramesRequired
	}
	if that.engine.fieldMap == nil {
		return ErrFieldFieldMapRequired
	}
	return nil
}

var (
	ErrFieldFramesRequired   = errors.New("frames required")
	ErrFieldFieldMapRequired = errors.New("field map required")
)
