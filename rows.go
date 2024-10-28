package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func GetRows(ctx context.Context, sql string, exec Executor, args ...any) (pgx.Rows, error) {
	return exec.Query(ctx, sql, args...)
}

type SQLRows func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error)

func NewSQLRows(sql string) SQLRows {
	return func(ctx context.Context, exec Executor, args ...any) (pgx.Rows, error) {
		return GetRows(ctx, sql, exec, args...)
	}
}
