package postgres_repos

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

func TestTagRepository(t *testing.T) {
	storage := CreatePostgresStorage(t)
	repo := &TagRepository{
		storage: storage,
	}

	abstract.TestTagRepositoryBehaviour(t, repo, GetCleanupFuncForTags(storage))
}

func GetCleanupFuncForTags(storage *Storage) func() error {
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
