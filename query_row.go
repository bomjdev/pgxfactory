package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type QueryRowFn func(ctx context.Context, exec Executor, args ...any) pgx.Row

func QueryRow(ctx context.Context, sql string, exec Executor, args ...any) pgx.Row {
	return exec.QueryRow(ctx, sql, args...)
}

func NewQueryRow(sql string) QueryRowFn {
	return func(ctx context.Context, exec Executor, args ...any) pgx.Row {
		return QueryRow(ctx, sql, exec, args...)
	}
}

type RowScanner[T any] interface {
	ScanRow(row pgx.Row) (T, error)
}

type RowScannerFn[T any] func(scan func(dest ...any) error) (T, error)

func NewRowScanner[T any](fn QueryRowFn, scanner RowScanner[T]) ExecScanFn[T] {
	return func(ctx context.Context, exec Executor, args ...any) (T, error) {
		return scanner.ScanRow(fn(ctx, exec, args...))
	}
}

func NewRowScannerFn[T any](fn QueryRowFn, scanner RowScannerFn[T]) ExecScanFn[T] {
	return func(ctx context.Context, exec Executor, args ...any) (T, error) {
		return scanner(fn(ctx, exec, args...).Scan)
	}
}
