package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Collector[T any] struct {
	many  ExecScanFn[[]T]
	first ExecScanFn[T]
	one   ExecScanFn[T]
}

func NewCollector[T any](getRows GetRowsFn, scanRow pgx.RowToFunc[T]) Collector[T] {
	return Collector[T]{
		many:  CollectRowsFactory(getRows, scanRow),
		first: CollectOneRowFactory(getRows, scanRow),
		one:   CollectExactlyOneRowFactory(getRows, scanRow),
	}
}

func (c Collector[T]) Many(ctx context.Context, exec Executor, args ...any) ([]T, error) {
	return c.many(ctx, exec, args...)
}

func (c Collector[T]) First(ctx context.Context, exec Executor, args ...any) (T, error) {
	return c.first(ctx, exec, args...)
}

func (c Collector[T]) One(ctx context.Context, exec Executor, args ...any) (T, error) {
	return c.one(ctx, exec, args...)
}
