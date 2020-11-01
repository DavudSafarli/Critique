package repositories

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/domain_test"
	"github.com/DavudSafarli/Critique/pkg/util"
)

func TestTagRepository(t *testing.T) {
	storage := util.CreatePostgresStorage(t)
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
