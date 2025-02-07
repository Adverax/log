package rex

import (
	"errors"
	"regexp"
)

type FrameBuilder struct {
	frame   *Entry
	pattern string
}

func NewFrameBuilder() *FrameBuilder {
	return &FrameBuilder{
		frame: &Entry{},
	}
}

func (that *FrameBuilder) WithPattern(pattern string) *FrameBuilder {
	that.pattern = pattern
	return that
}

func (that *FrameBuilder) WithGuard(guard Guard) *FrameBuilder {
	that.frame.guard = guard
	return that
}

func (that *FrameBuilder) Build() (*Entry, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	re, err := regexp.Compile(that.pattern)
	if err != nil {
		return nil, err
	}

	that.frame.re = re
	return that.frame, nil
}

func (that *FrameBuilder) checkRequiredFields() error {
	if that.pattern == "" {
		return ErrPatternRequired
	}
	return nil
}

var (
	ErrPatternRequired = errors.New("pattern is required")
)
