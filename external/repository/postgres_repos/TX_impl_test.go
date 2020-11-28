package postgres_repos

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestTx(t *testing.T) {
	storage, err := NewPostgresStorage(testing_utils.GetTestDbConnStr())
	require.Nil(t, err, "NewSingletonDbConnection Should not return error")

	t.Run("GetDB should return TX interface", func(t *testing.T) {
		ctx := context.Background()
		ctx, err = storage.BeginTx(ctx)
		require.Nil(t, err, "BeginTx Should not return error")

		db := storage.getDB(ctx)

		_, ok := db.(pgx.Tx)
		require.True(t, ok, "GetDB should return intance of TX interface, because we called BeginTx")
	})

	t.Run("GetDB should return pgx.Pool instance", func(t *testing.T) {
		ctx := context.Background()
		db := storage.getDB(ctx)

		require.Equal(t, storage.DB, db, "GetDB should return storage.DB, because we didn't call BeginTx")
	})
}
