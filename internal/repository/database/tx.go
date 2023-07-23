package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WithTxFunc func(ctx context.Context, tx *sqlx.Tx) error

func WithTx(ctx context.Context, db *sqlx.DB, fn WithTxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "repository.WithTx.BeginTxx")
	}

	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, "repository.WithTX.Rollback")
		}
		return errors.Wrap(err, "repository.WithTx.fn")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "repository.WithTx.Commit")
	}

	return nil
}
