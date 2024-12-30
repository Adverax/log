package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

type Renderer interface {
	Render(ctx context.Context, entry *Entry)
}

type BaseRenderer struct {
	formatter Formatter
	out       io.Writer
}

func NewRenderer(
	formatter Formatter,
	out io.Writer,
) *BaseRenderer {
	return &BaseRenderer{
		formatter: formatter,
		out:       out,
	}
}

func (that *BaseRenderer) Render(ctx context.Context, entry *Entry) {
	buffer := entry.Logger.GetBuffer()
	defer func() {
		entry.Buffer = nil
		buffer.Reset()
		entry.Logger.FreeBuffer(buffer)
	}()
	buffer.Reset()
	entry.Buffer = buffer

	that.render(entry)

	entry.Buffer = nil
}

func (that *BaseRenderer) render(entry *Entry) {
	serialized, err := that.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	if _, err := that.out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

type dummyRenderer struct{}

func (that *dummyRenderer) Render(ctx context.Context, entry *Entry) {
	// nothing
}
