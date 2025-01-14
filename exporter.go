package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

type Exporter interface {
	Export(ctx context.Context, entry *Entry)
}

type BaseExporter struct {
	formatter Formatter
	out       io.Writer
}

func NewRenderer(
	formatter Formatter,
	out io.Writer,
) *BaseExporter {
	return &BaseExporter{
		formatter: formatter,
		out:       out,
	}
}

func (that *BaseExporter) Export(ctx context.Context, entry *Entry) {
	buffer := entry.Logger.GetBuffer()
	defer func() {
		entry.Buffer = nil
		buffer.Reset()
		entry.Logger.FreeBuffer(buffer)
	}()
	buffer.Reset()
	entry.Buffer = buffer

	that.export(entry)

	entry.Buffer = nil
}

func (that *BaseExporter) export(entry *Entry) {
	serialized, err := that.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	if _, err := that.out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

type dummyExporter struct{}

func (that *dummyExporter) Export(ctx context.Context, entry *Entry) {
	// nothing
}
