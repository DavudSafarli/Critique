package impl

import (
	"github.com/DavudSafarli/Critique/domain/contracts"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

func TestAttchUc(t *testing.T) {
	t.Parallel()
	spec := testcase.NewSpec(t)
	spec.Parallel()
	SetupUsecaseDependencies(spec)
	getCtx := contracts.GetTxContextForTest

	spec.Describe(`AttachmentUsecases#CreateMany`, func(s *testcase.Spec) {
		attchsVar := testcase.Var{Name: `attchsVar`}
		getAttchs := func(t *testcase.T) []models.Attachment { return attchsVar.Get(t).([]models.Attachment) }
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			i := contracts.GetFeedbackID(t)
			return GetAttachmentUsecaseForTest(t).CreateAttachments(getCtx(t), getAttchs(t), i)
		}

		s.When(`Empty list of attachments given`, func(s *testcase.Spec) {
			attchsVar.Let(s, func(t *testcase.T) interface{} {
				return []models.Attachment{}
			})
			contracts.FeedbackID.LetValue(s, uint(99999)) // and whatever feedbackID is given
			s.Then(`It will return error`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Error(t, err)
				require.Len(t, attchs, 0)
			})
		})
		s.When(`Non-empty list of attachments given with ID of existing Feedback`, func(s *testcase.Spec) {
			attchsVar.Let(s, func(t *testcase.T) interface{} {
				return testing_utils.ExampleAttchSlice(4)
			})
			contracts.FeedbackID.Let(s, nil) // <-- bind Init func of feedbackID variable
			s.Then(`It should be alright`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err)
				require.Len(t, attchs, 4)
				// TODO: may be check if they have unique IDS?
				// TODO: I don't wanna add GetAttchs func to UC, or use its dep.s to check storage
			})
		})
		s.When(`Non-empty list of attachments given with ID of NON-existing Feedback`, func(s *testcase.Spec) {
			attchsVar.Let(s, func(t *testcase.T) interface{} {
				return testing_utils.ExampleAttchSlice(4)
			})
			contracts.FeedbackID.LetValue(s, uint(7373573)) // <--
			s.Then(`It will return error`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
			})
		})
	})

}
