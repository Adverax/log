package rex

import (
	"regexp"
)

type Frame interface {
	Parse(data []byte, fields map[string]string) (int, error)
}

type Guard interface {
	IsSatisfied(fields map[string]string) bool
}

type GuardFunc func(fields map[string]string) bool

func (that GuardFunc) IsSatisfied(fields map[string]string) bool {
	return that(fields)
}

type Entry struct {
	guard Guard
	re    *regexp.Regexp
}

func (that *Entry) Parse(data []byte, fields map[string]string) (int, error) {
	if that.guard != nil && !that.guard.IsSatisfied(fields) {
		return 0, nil
	}

	matches := that.re.FindSubmatch(data)
	if matches == nil {
		return 0, nil
	}

	keys := that.re.SubexpNames()

	for i, key := range keys {
		if key == "" {
			continue
		}
		fields[key] = string(matches[i])
	}

	return len(matches[0]), nil
}
