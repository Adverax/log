package fileExporter

import (
	"context"
	"fmt"
	"github.com/adverax/log"
	"io"
	"os"
)

type Exporter struct {
	formatter log.Formatter
	out       io.Writer
}

func New(
	formatter log.Formatter,
	out io.Writer,
) *Exporter {
	if out == nil {
		out = os.Stdout
	}

	return &Exporter{
		formatter: formatter,
		out:       out,
	}
}

func (that *Exporter) Export(ctx context.Context, entry *log.Entry) {
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

func (that *Exporter) export(entry *log.Entry) {
	serialized, err := that.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}

	if _, err := that.out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}
