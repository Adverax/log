package rex

import (
	"errors"
	"github.com/adverax/core"
	"regexp"
)

type FrameBuilder struct {
	*core.Builder
	frame   *Entry
	pattern string
}

func NewFrameBuilder() *FrameBuilder {
	return &FrameBuilder{
		Builder: core.NewBuilder("frame"),
		frame:   &Entry{},
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
	that.RequiredField(that.pattern, ErrPatternRequired)
	return that.ResError()
}

var (
	ErrPatternRequired = errors.New("pattern is required")
)
