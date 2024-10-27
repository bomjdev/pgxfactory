package pgxfactory

import (
	"github.com/jackc/pgx/v5"
	"testing"
)

func TestQueryRow(t *testing.T) {
	exec := NewQueryRow("SELECT 1")
	execScan := NewRowScannerFn(exec, func(scan func(dest ...any) error) (uint, error) {
		var x uint
		return x, scan(&x)
	})
	i, err := execScan(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatal(i)
	}

	execScan2 := NewRowScanner[testUintRowScanner](exec, testUintRowScanner{})
	s, err := execScan2(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}
	if s.x != 1 {
		t.Fatal(s.x)
	}
}

type testUintRowScanner struct {
	x uint
}

func (t testUintRowScanner) ScanRow(row pgx.Row) (testUintRowScanner, error) {
	return t, row.Scan(&t.x)
}
