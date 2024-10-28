package pgxfactory

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func Exec(ctx context.Context, sql string, exec Executor, args ...any) (pgconn.CommandTag, error) {
	return exec.Exec(ctx, sql, args...)
}

type ExecFn func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error)

func NewExec(sql string) ExecFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
		return Exec(ctx, sql, exec, args...)
	}
}

func (fn ExecFn) WithMiddleware(middlewares ...func(ExecFn) ExecFn) ExecFn {
	ret := fn
	for _, mw := range middlewares {
		ret = mw(ret)
	}
	return ret
}

func IsSelect(fn ExecFn) ExecFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
		tag, err := fn(ctx, exec, args...)
		if err != nil {
			return tag, err
		}
		if !tag.Select() {
			return tag, fmt.Errorf("expected select query, but got %s", tag.String())
		}
		return tag, nil
	}
}

func IsInsert(fn ExecFn) ExecFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
		tag, err := fn(ctx, exec, args...)
		if err != nil {
			return tag, err
		}
		if !tag.Insert() {
			return tag, fmt.Errorf("expected insert query, but got %s", tag.String())
		}
		return tag, nil
	}
}

func IsUpdate(fn ExecFn) ExecFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
		tag, err := fn(ctx, exec, args...)
		if err != nil {
			return tag, err
		}
		if !tag.Update() {
			return tag, fmt.Errorf("expected update query, but got %s", tag.String())
		}
		return tag, nil
	}
}

func IsDelete(fn ExecFn) ExecFn {
	return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
		tag, err := fn(ctx, exec, args...)
		if err != nil {
			return tag, err
		}
		if !tag.Delete() {
			return tag, fmt.Errorf("expected delete query, but got %s", tag.String())
		}
		return tag, nil
	}
}

func RowsAffected(n int64) func(ExecFn) ExecFn {
	return func(fn ExecFn) ExecFn {
		return func(ctx context.Context, exec Executor, args ...any) (pgconn.CommandTag, error) {
			tag, err := fn(ctx, exec, args...)
			if err != nil {
				return tag, err
			}
			if x := tag.RowsAffected(); x != n {
				return tag, fmt.Errorf("expected %d rows affected, but got %d", n, x)
			}
			return tag, nil
		}
	}
}
