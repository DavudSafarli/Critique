package postgres_repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgTxCtxKey struct{}

func (T *Storage) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := T.getDB(ctx).Begin(ctx) // tx -> Querier
	if err != nil {
		//fmt.Println("TX FAILED for context", i, err)
		return ctx, err
	}

	a := context.WithValue(ctx, PgTxCtxKey{}, tx)
	return a, nil
}

func (T *Storage) CommitTx(ctx context.Context) error {
	tx, err := T.lookupTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (T *Storage) RollbackTx(ctx context.Context) error {
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

// TXQuerier combines TX and Querier interfaces
type TXQuerier interface {
	pgxtype.Querier
	Begin(ctx context.Context) (pgx.Tx, error)
}

func (T *Storage) getDB(ctx context.Context) TXQuerier {
	tx, err := T.lookupTx(ctx)
	if err == nil {
		return tx
	}
	return T.DB
}

func (T *Storage) close(con TXQuerier) {
	if _, ok := con.(pgx.Tx); ok {
		return
	}
	if db, ok := con.(*pgxpool.Pool); ok {
		db.Close()
	}

}
