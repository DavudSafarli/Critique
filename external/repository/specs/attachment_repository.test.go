package specs

import (
	"context"

	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/spec_helper"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/testing_utils"
)

type AttchRequiredFuncs interface {
	GetAllAttachments() ([]models.Attachment, error)
}

type AttachmentRepositoryTester interface {
	abstract.AttachmentRepository
	GetAll(ctx context.Context) ([]models.Attachment, error)
}

// TestAttachmentRepositoryBehfeedbackIDaviour does what its name says
func TestAttachmentRepositoryBehaviour(s *testcase.Spec) {
	getCtx := spec_helper.GetTxContextForTest
	feedbackID := spec_helper.FeedbackID
	// creates  new feedback and returns its id
	s.Describe(`AttachmentRepository#CreateMany`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			slice := testing_utils.ExampleAttchSlice(3)
			return slice, spec_helper.GetAttachmentRepoForTest(t).CreateMany(getCtx(t), slice, feedbackID.Get(t).(uint))
		}
		s.When(`There is feedback beforehand`, func(s *testcase.Spec) {
			feedbackID.Let(s, nil) // <-- bind to the scope
			s.Then(`It should not return error and be able to retrieve`, func(t *testcase.T) {
				inserteds, err := subject(t)
				require.Nil(t, err)
				retrieveds, err := spec_helper.GetAttachmentRepoForTest(t).(AttachmentRepositoryTester).GetAll(getCtx(t))
				require.Nil(t, err)
				require.Equal(t, inserteds, retrieveds)
			})
		})
		s.When(`There is no feedback`, func(s *testcase.Spec) {
			feedbackID.LetValue(s, uint(999))
			s.Then(`It should return error, because creating attachments is impossible`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
			})
		})
	})

	s.Describe(`AttachmentRepository#GetByFeedbackID`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			return spec_helper.GetAttachmentRepoForTest(t).GetByFeedbackID(getCtx(t), feedbackID.Get(t).(uint))
		}
		s.When(`There are multiple feedbacks and attachments beforehand`, func(s *testcase.Spec) {
			feedbackID.Let(s, nil) // <-- bind it to scope
			s.Before(func(t *testcase.T) {
				i := feedbackID.Get(t).(uint)
				require.Nil(t, spec_helper.GetAttachmentRepoForTest(t).CreateMany(getCtx(t), testing_utils.ExampleAttchSlice(2), i))

				feedback2 := testing_utils.ExampleFeedback()
				require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), feedback2))
				require.Nil(t, spec_helper.GetAttachmentRepoForTest(t).CreateMany(getCtx(t), testing_utils.ExampleAttchSlice(3), feedback2.ID))
			})
			s.Then(`It should find them`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err, "No error pls")
				require.Len(t, attchs, 2)
			})
		})
		s.When(`There is feedbacks and but no attachments`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				f := testing_utils.ExampleFeedback()
				require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), f))
				feedbackID.Set(t, f.ID)
			})
			s.Then(`should return empty slice without error`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err, "No error pls")
				require.Len(t, attchs, 0)
			})
		})
		s.When(`There is no feedback`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				feedbackID.Set(t, uint(0))
			})
			s.Then(`should return empty slice without error`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err)
				require.Len(t, attchs, 0)
			})
		})
	})
}
