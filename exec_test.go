package pgxfactory

import (
	"testing"
)

func TestExec(t *testing.T) {
	execSQL := NewExec("SELECT 1")

	type testCase struct {
		name  string
		fn    ExecFn
		isErr bool
	}

	testCases := []testCase{
		{
			name:  "raw exec",
			fn:    execSQL,
			isErr: false,
		},
		{
			name:  "ensure select",
			fn:    execSQL.WithMiddleware(IsSelect),
			isErr: false,
		},
		{
			name:  "ensure insert",
			fn:    execSQL.WithMiddleware(IsInsert),
			isErr: true,
		},
		{
			name:  "ensure delete",
			fn:    execSQL.WithMiddleware(IsDelete),
			isErr: true,
		},
		{
			name:  "ensure update",
			fn:    execSQL.WithMiddleware(IsUpdate),
			isErr: true,
		},
		{
			name:  "ensure 1 rows affected",
			fn:    execSQL.WithMiddleware(RowsAffected(1)),
			isErr: false,
		},
		{
			name:  "ensure 0 rows affected",
			fn:    execSQL.WithMiddleware(RowsAffected(0)),
			isErr: true,
		},
		{
			name:  "ensure 2 rows affected",
			fn:    execSQL.WithMiddleware(RowsAffected(2)),
			isErr: true,
		},
		{
			name:  "ensure select 1 rows affected",
			fn:    execSQL.WithMiddleware(IsSelect, RowsAffected(1)),
			isErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.fn(ctx, pool)
			if (err != nil) != tc.isErr {
				t.Errorf("expected no error, but got %v", err)
			}
		})
	}

}
