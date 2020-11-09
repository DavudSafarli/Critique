package postgres_repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type TXImpl struct {
	storage *Storage
}
type PgTxCtxKey struct{}

func (T Storage) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := T.DB.Begin(ctx) // tx -> Querier
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, PgTxCtxKey{}, tx), nil
}

func (T Storage) CommitTx(ctx context.Context) error {
	tx, err := T.lookupTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (T Storage) RollbackTx(ctx context.Context) error {
	tx, err := T.lookupTx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback(ctx)
}

func (*Storage) lookupTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(PgTxCtxKey{}).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf(`no postgres tx in the given context`)
	}
	return tx, nil
}
