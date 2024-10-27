package pgxfactory

import (
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestSQL(t *testing.T) {
	sql := NewSQL("SELECT 1")
	execSQL := sql.exec.WithMiddleware(IsSelect, ExpectRowsAffected(1))

	_, err := execSQL(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}

	rowSQLVar := NewRowScannerFn(sql.row, func(scan func(dest ...any) error) (uint, error) {
		var x uint
		return x, scan(&x)
	})
	i, err := rowSQLVar(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatal(i)
	}

	rowsSQLVar := NewCollector[uint](sql.rows, pgx.RowTo[uint])
	rows, err := rowsSQLVar.many(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatal(len(rows))
	}
	if rows[0] != 1 {
		t.Fatal(rows[0])
	}
	row, err := rowsSQLVar.one(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if row != 1 {
		t.Fatal(row)
	}
	row, err = rowsSQLVar.first(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if row != 1 {
		t.Fatal(row)
	}
}
