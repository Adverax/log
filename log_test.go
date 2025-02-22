package log

import (
	"context"
	"os"
	"time"
)

func ExampleLog() {
	formatter, err := NewFormatterJsonBuilder().Build()
	if err != nil {
		panic(err)
	}

	logger, err := NewLogBuilder().
		WithLevel(InfoLevel).
		WithExporter(NewBaseExporter(formatter, os.Stdout)).
		WithHook(HookFunc(func(ctx context.Context, entry *Entry) error {
			entry.Time = time.Time{}
			return nil
		})).
		Build()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	logger.Info(ctx, "Hello, World!")

	// Output:
	// {"level":"info","msg":"Hello, World!","time":"0001-01-01 00:00:00"}
}
