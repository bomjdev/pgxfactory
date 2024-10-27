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
