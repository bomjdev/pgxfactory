package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type CollectFn[T, R any] func(rows pgx.Rows, fn pgx.RowToFunc[T]) (R, error)

func CollectFactory[T, R any](getRows GetRowsFn, collect CollectFn[T, R], scanRow pgx.RowToFunc[T]) ExecScanFn[R] {
	return func(ctx context.Context, exec Executor, args ...any) (R, error) {
		rows, err := getRows(ctx, exec, args...)
		if err != nil {
			return *new(R), err
		}
		return collect(rows, scanRow)
	}
}

func CollectRowsFactory[T any](getRows GetRowsFn, scanRow pgx.RowToFunc[T]) ExecScanFn[[]T] {
	return CollectFactory(getRows, pgx.CollectRows[T], scanRow)
}

func CollectExactlyOneRowFactory[T any](getRows GetRowsFn, scanRow pgx.RowToFunc[T]) ExecScanFn[T] {
	return CollectFactory(getRows, pgx.CollectExactlyOneRow[T], scanRow)
}

func CollectOneRowFactory[T any](getRows GetRowsFn, scanRow pgx.RowToFunc[T]) ExecScanFn[T] {
	return CollectFactory(getRows, pgx.CollectOneRow[T], scanRow)
}
