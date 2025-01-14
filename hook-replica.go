package log

import (
	"context"
)

type HookReplica struct {
	exporter Exporter
}

func NewHookReplica(
	exporter Exporter,
) *HookReplica {
	return &HookReplica{
		exporter: exporter,
	}
}

func (that *HookReplica) Fire(ctx context.Context, entry *Entry) error {
	that.exporter.Export(ctx, entry)
	return nil
}
