package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Query[T any] func(ctx context.Context, exec DB, args ...any) (T, error)

type RowsCollector[T, R any] func(rows pgx.Rows, rowTo pgx.RowToFunc[T]) (R, error)

func NewQuery[T, R any](getRows SQLRows, collect RowsCollector[T, R], rowTo pgx.RowToFunc[T]) Query[R] {
	return func(ctx context.Context, exec DB, args ...any) (R, error) {
		var v R
		rows, err := getRows(ctx, exec, args...)
		if err != nil {
			return v, err
		}
		return collect(rows, rowTo)
	}
}

func NewQueryAll[T any](getRows SQLRows, rowTo pgx.RowToFunc[T]) Query[[]T] {
	return NewQuery(getRows, pgx.CollectRows[T], rowTo)
}

func NewQueryOne[T any](getRows SQLRows, rowTo pgx.RowToFunc[T]) Query[T] {
	return NewQuery(getRows, pgx.CollectExactlyOneRow[T], rowTo)
}

func NewQueryFirst[T any](getRows SQLRows, rowTo pgx.RowToFunc[T]) Query[T] {
	return NewQuery(getRows, pgx.CollectOneRow[T], rowTo)
}

type Querier[T any] struct {
	many  Query[[]T]
	first Query[T]
	one   Query[T]
}

func NewQuerier[T any](queryFn SQLRows, rowTo pgx.RowToFunc[T]) Querier[T] {
	return Querier[T]{
		many:  NewQueryAll(queryFn, rowTo),
		first: NewQueryFirst(queryFn, rowTo),
		one:   NewQueryOne(queryFn, rowTo),
	}
}

func (s Querier[T]) All(ctx context.Context, exec DB, args ...any) ([]T, error) {
	return s.many(ctx, exec, args...)
}

func (s Querier[T]) First(ctx context.Context, exec DB, args ...any) (T, error) {
	return s.first(ctx, exec, args...)
}

func (s Querier[T]) One(ctx context.Context, exec DB, args ...any) (T, error) {
	return s.one(ctx, exec, args...)
}
