package pgxfactory

import (
	"testing"
)

func TestExec(t *testing.T) {
	fn := NewExecFn("SELECT 1")
	fn = fn.WithMiddleware(IsSelect, ExpectRowsAffected(1))
	tag, err := fn(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tag)
}
