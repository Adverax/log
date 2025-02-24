package logFileRotator

import (
	"fmt"
	"os"
	"path/filepath"
)

type Builder struct {
	engine *Engine
}

func NewBuilder() *Builder {
	return &Builder{
		engine: &Engine{
			options: Options{
				maxSize:    10000000,
				maxAge:     30,
				maxBackups: 30,
				localTime:  false,
				timeFormat: "2006-01-02T15-04-05.000",
			},
		},
	}
}

func (that *Builder) WithFileName(filename string) *Builder {
	that.engine.options.fileName = filename
	return that
}

func (that *Builder) WithMaxSize(maxSize int) *Builder {
	that.engine.options.maxSize = maxSize
	return that
}

func (that *Builder) WithMaxAge(maxAge int) *Builder {
	that.engine.options.maxAge = maxAge
	return that
}

func (that *Builder) WithMaxBackups(maxBackups int) *Builder {
	that.engine.options.maxBackups = maxBackups
	return that
}

func (that *Builder) WithLocalTime(localTime bool) *Builder {
	that.engine.options.localTime = localTime
	return that
}

func (that *Builder) Build() (*Engine, error) {
	if err := that.updateDefaultFields(); err != nil {
		return nil, err
	}

	return that.engine, nil
}

func (that *Builder) updateDefaultFields() error {
	if that.engine.options.fileName == "" {
		var err error
		that.engine.options.fileName, err = that.createDefaultFileName()
		if err != nil {
			return err
		}
	}

	return nil
}

func (that *Builder) createDefaultFileName() (string, error) {
	dir, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "log"))
	if err != nil {
		return "", fmt.Errorf("Error get abs path:%v", err)
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Error create log dir:%v", err)
	}

	file := filepath.Base(os.Args[0]) + ".log"
	return filepath.Join(dir, file), nil
}
