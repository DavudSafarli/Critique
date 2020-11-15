package postgres_repos

import (
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"testing"

	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

func TestFeedbackRepository(t *testing.T) {
	storage, err := NewSingletonDbConnection(testing_utils.GetTestDbConnStr())
	if err != nil {
		panic(err)
	}
	repo := &FeedbackRepository{storage}

	abstract.TestFeedbackRepositoryBehaviour(t, repo, testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
}
