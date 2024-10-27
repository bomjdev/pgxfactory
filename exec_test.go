package pgxfactory

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

func TestExec(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	fn := NewExecFn("SELECT 1")
	fn = fn.WithMiddleware(IsSelect, ExpectRowsAffected(1))
	tag, err := fn(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tag)
}
