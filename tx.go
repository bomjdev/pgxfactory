package pgxfactory

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func RunInTransaction(ctx context.Context, beginner Beginner, fn func(tx pgx.Tx) error) error {
	tx, err := beginner.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}

	defer func() { CommitOrRollback(ctx, tx, err) }()

	if err = fn(tx); err != nil {
		return fmt.Errorf("run in transaction: %w", err)
	}

	return nil
}

/*
CommitOrRollback commits or rollbacks the transaction depending on the error
Usage: defer func() { CommitOrRollback(ctx, tx, err) }()
*/
func CommitOrRollback(ctx context.Context, tx pgx.Tx, cause error) {
	if cause != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			_ = log.Output(3, fmt.Sprintf("rollback error: %s; rollback caused by: %s", err, cause))
		}
	}
	if err := tx.Commit(ctx); err != nil {
		_ = log.Output(3, fmt.Sprintf("commit error: %s", err))
	}
}
