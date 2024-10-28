package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func Query(ctx context.Context, sql string, exec Executor, args ...any) (pgx.Rows, error) {
	return exec.Query(ctx, sql, args...)
}

type QueryFn func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error)

func NewQuery(sql string) QueryFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error) {
		return Query(ctx, sql, exec, args...)
	}
}
