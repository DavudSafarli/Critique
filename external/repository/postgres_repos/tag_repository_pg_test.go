package postgres_repos

import (
	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"testing"
)

func TestTagRepository(t *testing.T) {
	storage, err := NewSingletonDbConnection(testing_utils.GetTestDbConnStr())
	if err != nil {
		panic(err)
	}
	repo := &TagRepository{storage}

	abstract.TestTagRepositoryBehaviour(t, repo, testing_utils.TruncateTestTables(t, "tags"))
}
