package postgres_repos

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

func TestFeedbackRepository(t *testing.T) {
	storage := vars.storage
	repo := &FeedbackRepository{storage}

	abstract.TestFeedbackRepositoryBehaviour(t, repo, GetCleanupFuncForFeedbacks(storage))
}

func GetCleanupFuncForFeedbacks(storage *Storage) func() error {
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
