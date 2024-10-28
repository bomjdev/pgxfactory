package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type ExecScanFn[T any] func(ctx context.Context, exec Executor, args ...any) (T, error)

type ScanFn[T, R any] func(rows pgx.Rows, fn pgx.RowToFunc[T]) (R, error)

func NewScanAll[T any](getRows QueryFn, scanRow pgx.RowToFunc[T]) ExecScanFn[[]T] {
	return NewScan(getRows, pgx.CollectRows[T], scanRow)
}

func NewScanOne[T any](getRows QueryFn, scanRow pgx.RowToFunc[T]) ExecScanFn[T] {
	return NewScan(getRows, pgx.CollectExactlyOneRow[T], scanRow)
}

func NewScanFirst[T any](getRows QueryFn, scanRow pgx.RowToFunc[T]) ExecScanFn[T] {
	return NewScan(getRows, pgx.CollectOneRow[T], scanRow)
}

func NewScan[T, R any](getRows QueryFn, collect ScanFn[T, R], scanRow pgx.RowToFunc[T]) ExecScanFn[R] {
	return func(ctx context.Context, exec Executor, args ...any) (R, error) {
		var v R
		rows, err := getRows(ctx, exec, args...)
		if err != nil {
			return v, err
		}
		return collect(rows, scanRow)
	}
}

type RowScanner[T any] struct {
	many  ExecScanFn[[]T]
	first ExecScanFn[T]
	one   ExecScanFn[T]
}

func NewScanner[T any](queryFn QueryFn, scanRow pgx.RowToFunc[T]) RowScanner[T] {
	return RowScanner[T]{
		many:  NewScanAll(queryFn, scanRow),
		first: NewScanFirst(queryFn, scanRow),
		one:   NewScanOne(queryFn, scanRow),
	}
}

func (s RowScanner[T]) All(ctx context.Context, exec Executor, args ...any) ([]T, error) {
	return s.many(ctx, exec, args...)
}

func (s RowScanner[T]) First(ctx context.Context, exec Executor, args ...any) (T, error) {
	return s.first(ctx, exec, args...)
}

func (s RowScanner[T]) One(ctx context.Context, exec Executor, args ...any) (T, error) {
	return s.one(ctx, exec, args...)
}
