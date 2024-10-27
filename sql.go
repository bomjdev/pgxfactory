package pgxfactory

func NewSQL(sql string) SQLExecutor {
	return SQLExecutor{
		exec: NewExecFn(sql),
		row:  NewQueryRow(sql),
		rows: NewGetRows(sql),
	}
}

type SQLExecutor struct {
	exec ExecFn
	row  QueryRowFn
	rows GetRowsFn
}
