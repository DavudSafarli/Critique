package postgres_repos

import (
	"testing"

	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/adamluzsi/testcase"
)

func TestPG(t *testing.T) {
	s := testcase.NewSpec(t)
	t.Parallel()
	s.Parallel()
	SetupPostgresStorageVar(s)
	s.Describe(`TestFeedbackRepositoryBehaviour`, func(s *testcase.Spec) {
		contracts.TestFeedbackRepositoryBehaviour(s)
	})
	s.Describe(`TestAttachmentRepositoryBehaviour`, func(s *testcase.Spec) {
		contracts.TestAttachmentRepositoryBehaviour(s)
	})
	s.Describe(`TestTagRepositoryBehaviour`, func(s *testcase.Spec) {
		contracts.TestTagRepositoryBehaviour(s)
	})
}
