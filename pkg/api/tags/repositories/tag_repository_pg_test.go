package repositories

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/domain_test"
	"github.com/DavudSafarli/Critique/pkg/util"
)

func CreateStorage(t *testing.T) *postgres.Storage {
	pgConnectionString := util.RunPostgresDockerAndGetConnectionString(t)
	storage, err := postgres.NewDbConnection(pgConnectionString)
	util.MigrateDatabase(t, pgConnectionString)
	if err != nil {
		t.Fatalf("Failed to create a new Storage: %s", err)
	}
	return storage
}

func TestRepository(t *testing.T) {
	storage := CreateStorage(t)
	repo := &TagRepository{
		storage: storage,
	}

	domain_test.TestTagRepositoryBehaviour(t, repo, GetCleanupFuncForTags(storage))
}

func GetCleanupFuncForTags(storage *postgres.Storage) func() error {
	return func() error {
		q := storage.SB.Delete("tags")

		sql, args, err := q.ToSql()
		if err != nil {
			return err
		}
		_, err = storage.DB.Exec(context.Background(), sql, args...)
		if err != nil {
			return err
		}
		return nil
	}
}
