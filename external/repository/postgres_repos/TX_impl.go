package postgres_repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgTxCtxKey struct{}

// txQuerier combines TX.Begin with Querier interface
type txQuerier interface {
	pgxtype.Querier
	Begin(ctx context.Context) (pgx.Tx, error)
}

func (T Storage) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := getDB(ctx, T.DB).Begin(ctx) // tx -> Querier
	if err != nil {
		return ctx, err
	}

	ctxWithTx := context.WithValue(ctx, PgTxCtxKey{}, tx)
	return ctxWithTx, nil
}

func lookupTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(PgTxCtxKey{}).(pgx.Tx)

	return tx, ok
}
func getDB(ctx context.Context, db *pgxpool.Pool) txQuerier {
	tx, ok := lookupTx(ctx)
	if ok {
		return tx
	}
	return db
}

func (T Storage) CommitTx(ctx context.Context) error {
	tx, ok := lookupTx(ctx)
	if !ok {
		return fmt.Errorf(`no postgres tx in the given context`)
	}
	return tx.Commit(ctx)
}

func (T Storage) RollbackTx(ctx context.Context) error {
	tx, ok := lookupTx(ctx)
	if !ok {
		return fmt.Errorf(`no postgres tx in the given context`)
	}
	return tx.Rollback(ctx)
}
