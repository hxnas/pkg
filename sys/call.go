package sys

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type Caller func(ctx context.Context) (err error)

func (f Caller) Call(ctx context.Context) (err error)  { return f(ctx) }
func (f Caller) Wrap(ctx context.Context) func() error { return func() error { return f(ctx) } }

// 并行执行
func ParallelCall(ctx context.Context, callers ...Caller) error {
	g, c := errgroup.WithContext(ctx)
	for _, caller := range callers {
		g.Go(caller.Wrap(c))
	}
	return g.Wait()
}

// 顺序执行
func Call(ctx context.Context, callers ...Caller) (err error) {
	for _, caller := range callers {
		if err = caller.Call(ctx); err != nil {
			break
		}
	}
	return
}

// 顺序执行
func SkipErrCall(ctx context.Context, callers ...Caller) {
	for _, caller := range callers {
		if err := caller.Call(ctx); err != nil {
			slog.WarnContext(ctx, "caller error", "err", err)
			break
		}
	}
}
