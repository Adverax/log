package rex

import "github.com/adverax/log"

type Engine struct {
	frames []Frame
}

func NewEngine() *Engine {
	return &Engine{
		frames: make([]Frame, 0, 8),
	}
}

func (that *Engine) AddFrame(frame ...Frame) {
	that.frames = append(that.frames, frame...)
}

func (that *Engine) Parse(data []byte, entry log.Entry) (int, error) {
	var nn int
	for _, frame := range that.frames {
		n, err := frame.Parse(data, entry)
		if err != nil {
			return n, err
		}
		nn += n
	}
	return nn, nil
}
