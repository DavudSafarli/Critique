package postgres_repos

import (
	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"testing"
)

func TestAttachmentRepository(t *testing.T) {
	storage, err := NewSingletonDbConnection(testing_utils.GetTestDbConnStr())
	if err != nil {
		panic(err)
	}
	attchrepo := AttachmentRepository{storage}
	fdbkRepo := FeedbackRepository{storage}
	abstract.TestAttachmentRepositoryBehaviour(t, attchrepo, fdbkRepo, testing_utils.TruncateTestTables(t, "attachments", "feedbacks"))
}
