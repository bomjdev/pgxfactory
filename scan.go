package pgxfactory

import "context"

type ExecScanFn[T any] func(ctx context.Context, exec Executor, args ...any) (T, error)
