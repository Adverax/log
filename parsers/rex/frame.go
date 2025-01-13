package rex

import (
	"github.com/adverax/log"
	"regexp"
)

type Frame interface {
	Parse(data []byte, entry *log.Entry) (int, error)
}

type Entry struct {
	re *regexp.Regexp
}

func NewFrame(pattern string) (*Entry, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Entry{re: re}, nil
}

func (that *Entry) Parse(data []byte, entry *log.Entry) (int, error) {
	matches := that.re.FindSubmatch(data)
	if matches == nil {
		return 0, nil
	}

	keys := that.re.SubexpNames()

	for i, key := range keys {
		if key == "" {
			continue
		}
		entry.Data[key] = string(matches[i])
	}

	return len(matches[0]), nil
}
