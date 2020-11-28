package repository

import (
	"testing"

	"github.com/DavudSafarli/Critique/external/repository/specs"
	"github.com/adamluzsi/testcase"
)

func TestAttachmentRepositoryBehaviour(t *testing.T) {
	spec := testcase.NewSpec(t)
	t.Parallel()
	spec.Parallel()
	specs.TestAttachmentRepositoryBehaviour(spec)
}
func TestFeedbackRepositoryBehaviour(t *testing.T) {
	spec := testcase.NewSpec(t)
	t.Parallel()
	spec.Parallel()
	specs.TestFeedbackRepositoryBehaviour(spec)
}
func TestTagRepositoryBehaviour(t *testing.T) {
	spec := testcase.NewSpec(t)
	t.Parallel()
	spec.Parallel()
	specs.TestTagRepositoryBehaviour(spec)
}
