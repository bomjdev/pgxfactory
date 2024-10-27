package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type GetRowsFn func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error)

func Query(ctx context.Context, sql string, exec Executor, args ...any) (pgx.Rows, error) {
	return exec.Query(ctx, sql, args...)
}

func NewGetRows(sql string) GetRowsFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error) {
		return Query(ctx, sql, exec, args...)
	}
}

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

func CollectScannerFactory[T, R any](getRows GetRowsFn, scanRow pgx.RowToFunc[T]) func(collect CollectFn[T, R]) ExecScanFn[R] {
	return func(collect CollectFn[T, R]) ExecScanFn[R] {
		return CollectFactory(getRows, collect, scanRow)
	}
}
