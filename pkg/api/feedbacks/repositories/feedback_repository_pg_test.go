package repositories

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/domain_test"
	"github.com/DavudSafarli/Critique/pkg/util"
)

func TestFeedbackRepository(t *testing.T) {
	storage := util.CreatePostgresStorage(t)
	repo := &FeedbackRepository{
		storage: storage,
	}

	domain_test.TestFeedbackRepositoryBehaviour(t, repo, GetCleanupFuncForFeedbacks(storage))
}

func GetCleanupFuncForFeedbacks(storage *postgres.Storage) func() error {
	return func() error {
		q := storage.SB.Delete("feedbacks")

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
