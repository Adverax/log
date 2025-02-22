package log

import (
	"context"
	"fmt"
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
	logger.
		WithFields(Fields{"key": "value"}).
		WithError(fmt.Errorf("invalid value")).
		Error(ctx, "Hello, World2!")

	// Output:
	// {"level":"info","msg":"Hello, World!","time":"0001-01-01 00:00:00"}
	// {"data":{"error":"invalid value","key":"value"},"level":"error","msg":"Hello, World2!","time":"0001-01-01 00:00:00"}
}
