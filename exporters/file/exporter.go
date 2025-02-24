package fileExporter

import (
	"context"
	"fmt"
	"github.com/adverax/log"
	"io"
	"os"
)

type Exporter interface {
	Export(ctx context.Context, entry *log.Entry)
}

type BaseExporter struct {
	formatter log.Formatter
	out       io.Writer
}

func New(
	formatter log.Formatter,
	out io.Writer,
) *BaseExporter {
	if out == nil {
		out = os.Stdout
	}

	return &BaseExporter{
		formatter: formatter,
		out:       out,
	}
}

func (that *BaseExporter) Export(ctx context.Context, entry *log.Entry) {
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

func (that *BaseExporter) export(entry *log.Entry) {
	serialized, err := that.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	if _, err := that.out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}
