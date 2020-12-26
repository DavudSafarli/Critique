package impl

import (
	"errors"
	"testing"

	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/domain/usecases"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository/mocks"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func LetValueUnsafe(v testcase.Var, s *testcase.Spec, T interface{}) {
	v.Let(s, func(t *testcase.T) interface{} {
		return T
	})
}

func TestFeedbackUc(t *testing.T) {
	t.Parallel()
	spec := testcase.NewSpec(t)
	spec.Parallel()
	SetupUsecaseDependencies(spec)

	getctx := contracts.GetTxContextForTest
	getFeedbackUc := GetFeedbackUsecaseForTest
	spec.Describe(`#CreateFeedback`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) (models.Feedback, error) {
			return getFeedbackUc(t).CreateFeedback(getctx(t), *contracts.GetFeedback(t))
		}
		s.When(`We provide invalid feedback`, func(s *testcase.Spec) {
			LetValueUnsafe(contracts.Feedback, s, &models.Feedback{}) // invalid feedback
			s.Then(`It will return an error`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
			})
		})
		s.When(`We provide valid feedback withOUT attachment`, func(s *testcase.Spec) {
			LetValueUnsafe(contracts.Feedback, s, testing_utils.ExampleFeedback()) // valid feedback
			s.Then(`It should be retrievable`, func(t *testcase.T) {
				f, err := subject(t)
				require.Nil(t, err)
				retrieved, err := getFeedbackUc(t).GetFeedbackDetails(getctx(t), f.ID)
				require.Nil(t, err)
				require.Equal(t, retrieved, f)
				require.Len(t, retrieved.Attachments, 0)
			})
		})
		s.When(`We provide valid feedback WITH attachments`, func(s *testcase.Spec) {
			f := testing_utils.ExampleFeedback()
			f.Attachments = testing_utils.ExampleAttchSlice(4)
			LetValueUnsafe(contracts.Feedback, s, f) // valid feedback
			s.Then(`It and its attachments should be retrievable`, func(t *testcase.T) {
				inserted, err := subject(t)
				require.Nil(t, err)
				retrieved, err := getFeedbackUc(t).GetFeedbackDetails(getctx(t), inserted.ID)
				require.Nil(t, err)
				require.Equal(t, retrieved, inserted)
			})
		})
		s.When(`Feedback and attachments are valid, but Attachment creation fails for some reason`, func(s *testcase.Spec) {
			f := testing_utils.ExampleFeedback()
			f.Attachments = testing_utils.ExampleAttchSlice(4)
			LetValueUnsafe(contracts.Feedback, s, f) // valid feedback

			s.Before(func(t *testcase.T) { // replace storage with mock
				mockRepo := mocks.FailingAttachmentFakeStorage{
					Storage: GetFeedbackUsecaseForTest(t).storage,
				}
				mockRepo.On("CreateManyAttachments", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Any error"))
				uc := GetFeedbackUsecaseForTest(t)
				uc.storage = mockRepo
				FeedbackUsecaseForTest.Set(t, uc)
			})
			s.Then(`It should maintain atomicity(Db should not have a change)`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
				allFeedbacks, err := getFeedbackUc(t).GetFeedbacksWithPagination(getctx(t), usecases.Pagination{Skip: 0, Limit: 100})
				require.Nil(t, err)
				require.Len(t, allFeedbacks, 0, "Should not create any feedback when Attachment creation failed")
			})
		})
	})
}
