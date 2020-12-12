package impl

import (
	"errors"
	"github.com/DavudSafarli/Critique/domain/contracts"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
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
				mockRepo := &mocks.FailingAttachmentFakeStorage{
					Storage: GetFeedbackUsecaseForTest(t).storage,
				}
				mockRepo.On("CreateManyAttachments", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Any error"))
				uc := GetFeedbackUsecaseForTest(t)
				uc.storage = mockRepo
			})
			s.Then(`It should maintain atomicity(Db should not have a change)`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
				allFeedbacks, err := getFeedbackUc(t).GetFeedbacksWithPagination(getctx(t), feedback_usecases.Pagination{Skip: 0, Limit: 100})
				require.Nil(t, err)
				require.Len(t, allFeedbacks, 0, "Should not create any feedback when Attachment creation failed")
			})
		})
	})
}

// old tests, don't look
//
// type initFuncReturnType func() (FeedbackUsecasesImpl, abstract.FeedbackRepository, abstract.AttachmentRepository)
//
//var getVars = (func() initFuncReturnType {
//	var uc FeedbackUsecasesImpl
//	var feedbackRepo abstract.FeedbackRepository
//	var attachmentRepo abstract.AttachmentRepository
//	var oncer sync.Once
//	return func() (FeedbackUsecasesImpl, abstract.FeedbackRepository, abstract.AttachmentRepository) {
//		oncer.Do(func() {
//			driver, connstr := "pg", testing_utils.GetTestDbConnStr()
//			feedbackRepo = repository.NewFeedbackRepository(driver, connstr)
//			attachmentRepo = repository.NewAttachmentRepository(driver, connstr)
//			uc = NewFeedbackUsecasesImpl(feedbackRepo, attachmentRepo).(FeedbackUsecasesImpl)
//		})
//		return uc, feedbackRepo, attachmentRepo
//	}
//})()
//
//func TestCreateFeedback(t *testing.T) {
//	t.Run("Should return error, Given unvalid Feedback", func(t *testing.T) {
//		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
//		usecase, _, _ := getVars()
//		invalidFeedback := models.Feedback{Title: ""}
//		_, err := usecase.CreateFeedback(context.Background(), invalidFeedback)
//		require.Error(t, err, "Should return error on invalid feedback")
//	})
//
//	// happy paths
//	t.Run("Should create Feedback and Attachments given a valid Feedback model", func(t *testing.T) {
//		// arrange
//		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
//		usecase, _, _ := getVars()
//		// act
//		validInput := models.Feedback{
//			Title:     "Non empty title",
//			Body:      "",
//			CreatedBy: "",
//			CreatedAt: uint(time.Now().Unix()),
//		}
//		f, err := usecase.CreateFeedback(context.Background(), validInput)
//		// assert
//		require.Nil(t, err, "Should not return error on valid input")
//		// assert models have assigned ID
//		require.NotZero(t, f.ID, "Feedback should have been assigned an ID")
//
//		storedFeedback, err := usecase.GetFeedbackDetails(context.Background(), f.ID)
//		require.Nil(t, err, "Should not return error on valid ID")
//		require.NotNil(t, storedFeedback, "Should not return nil Feedback on valid ID")
//		require.Equal(t, storedFeedback.Title, validInput.Title, "Should have the same title")
//	})
//
//	t.Run("Should create Feedback and Attachments given a valid Feedback model", func(t *testing.T) {
//		// arrange
//		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
//		usecase, _, _ := getVars()
//		// act
//		validInput := models.Feedback{
//			Title:     "Non empty title",
//			Body:      "",
//			CreatedBy: "",
//			CreatedAt: uint(time.Now().Unix()),
//			Attachments: []models.Attachment{
//				{Name: "bla bla2", Path: "/somewhere1"},
//				{Name: "bla bla2", Path: "/somewhere2"},
//			},
//		}
//		f, err := usecase.CreateFeedback(context.Background(), validInput)
//		// assert
//		require.Nil(t, err, "Should not return error on valid input")
//		// assert models have assigned ID
//		require.NotZero(t, f.ID, "Feedback should have been assigned an ID")
//		for _, attch := range f.Attachments {
//			require.NotZero(t, attch.ID, "Attachment should have been assigned an ID")
//		}
//
//		storedFeedback, err := usecase.GetFeedbackDetails(context.Background(), f.ID)
//		require.Nil(t, err, "Should not return error on valid ID")
//		require.NotNil(t, storedFeedback, "Should not return nil Feedback on valid ID")
//		require.Equal(t, storedFeedback.Title, validInput.Title, "Should have the same title")
//		require.Equal(t, len(validInput.Attachments), len(storedFeedback.Attachments), "Should have stored all attachments/fetched all attachments")
//		for i := 0; i < len(validInput.Attachments); i++ {
//			storedAttch := storedFeedback.Attachments[i]
//			insertedAttch := validInput.Attachments[i]
//			require.Equal(t, insertedAttch.Name, storedAttch.Name, "Stored Attachment should have the same Name")
//			require.Equal(t, insertedAttch.Path, storedAttch.Path, "Stored Attachment should have the same Path")
//		}
//	})
//	//
//	//t.Run("Should maintain atomcity. If Attachments can't be created, database should not have a change", func(t *testing.T) {
//	//	// arrange
//	//	t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
//	//	driver, connstr := "pg", testing_utils.GetTestDbConnStr()
//	//	feedbackRepo := repository.NewFeedbackRepository(driver, connstr)
//	//	attachmentRepo := mocks.MockAttachmentRepository{}
//	//
//	//	attachmentRepo.On("CreateMany", mock.Anything, mock.Anything, mock.Anything).Return([]models.Attachment{}, errors.New("Any error"))
//	//	usecase := NewFeedbackUsecasesImpl(feedbackRepo, attachmentRepo).(FeedbackUsecasesImpl)
//	//
//	//	//act
//	//	validInput := models.Feedback{
//	//		Title:     "Non empty title",
//	//		Body:      "",
//	//		CreatedBy: "",
//	//		CreatedAt: uint(time.Now().Unix()),
//	//		Attachments: []models.Attachment{
//	//			{Name: "bla bla2", Path: "/somewhere1"},
//	//			{Name: "bla bla2", Path: "/somewhere2"},
//	//		},
//	//	}
//	//
//	//	_, err := usecase.CreateFeedback(context.Background(), validInput)
//	//	// assert
//	//	require.Error(t, err, "Should return error")
//	//	fdbcks, err := usecase.GetFeedbacksWithPagination(context.Background(), feedback_usecases.Pagination{Skip: 0, Limit: 100})
//	//	require.Nil(t, err, "No error pls")
//	//	require.Zero(t, len(fdbcks), "Should not create any feedback when Attachment creation failed")
//	//
//	//})
//
//}
//
//func TestGetFeedbacksWithPagination(t *testing.T) {
//	t.Run("Doesn't use FeedbackRepository and 0-length slice and error when pagination limit is 0", func(t *testing.T) {
//		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
//		// arrange
//		usecase, _, _ := getVars()
//
//		// act
//		pagination := feedback_usecases.Pagination{Limit: 0}
//		feedbacks, err := usecase.GetFeedbacksWithPagination(context.Background(), pagination)
//
//		// assert
//		assert.Equal(t, ZeroLimitPaginationErr, err, "Should return non-nil error when limit is 0")
//		assert.Equal(t, 0, len(feedbacks), "Should return nil(0-length) slice when pagination limit is 0")
//	})
//}
