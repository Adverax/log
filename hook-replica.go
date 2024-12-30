package log

import (
	"context"
)

type HookReplica struct {
	renderer Renderer
}

func NewHookReplica(
	renderer Renderer,
) *HookReplica {
	return &HookReplica{
		renderer: renderer,
	}
}

func (that *HookReplica) Fire(ctx context.Context, entry *Entry) error {
	that.renderer.Render(ctx, entry)
	return nil
}
