package impl

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
	"github.com/DavudSafarli/Critique/external/repository"
	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type initFuncReturnType func() (FeedbackUsecasesImpl, abstract.FeedbackRepository, abstract.AttachmentRepository)

var getVars = (func() initFuncReturnType {
	var uc FeedbackUsecasesImpl
	var feedbackRepo abstract.FeedbackRepository
	var attachmentRepo abstract.AttachmentRepository
	var oncer sync.Once
	return func() (FeedbackUsecasesImpl, abstract.FeedbackRepository, abstract.AttachmentRepository) {
		oncer.Do(func() {
			driver, connstr := "pg", testing_utils.GetTestDbConnStr()
			feedbackRepo = repository.NewFeedbackRepository(driver, connstr)
			attachmentRepo = repository.NewAttachmentRepository(driver, connstr)
			uc = NewFeedbackUsecasesImpl(feedbackRepo, attachmentRepo).(FeedbackUsecasesImpl)
		})
		return uc, feedbackRepo, attachmentRepo
	}
})()

func TestCreateFeedback(t *testing.T) {
	t.Run("Should return error, Given unvalid Feedback", func(t *testing.T) {
		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
		usecase, _, _ := getVars()
		invalidFeedback := models.Feedback{Title: ""}
		_, err := usecase.CreateFeedback(context.Background(), invalidFeedback)
		require.Error(t, err, "Should return error on invalid feedback")
	})

	// happy path
	t.Run("Should create Feedback and Attachments given a valid Feedback model", func(t *testing.T) {
		// arrange
		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
		usecase, _, _ := getVars()
		// act
		validInput := models.Feedback{
			Title:     "Non empty title",
			Body:      "",
			CreatedBy: "",
			CreatedAt: uint(time.Now().Unix()),
			Attachments: []models.Attachment{
				{Name: "bla bla2", Path: "/somewhere1"},
				{Name: "bla bla2", Path: "/somewhere2"},
			},
		}
		f, err := usecase.CreateFeedback(context.Background(), validInput)
		// assert
		require.Nil(t, err, "Should not return error on valid input")
		// assert models have assigned ID
		require.NotZero(t, f.ID, "Feedback should have been assigned an ID")
		for _, attch := range f.Attachments {
			require.NotZero(t, attch.ID, "Attachment should have been assigned an ID")
		}

		storedFeedback, err := usecase.GetFeedbackDetails(context.Background(), f.ID)
		require.Nil(t, err, "Should not return error on valid ID")
		require.NotNil(t, storedFeedback, "Should not return nil Feedback on valid ID")
		require.Equal(t, storedFeedback.Title, validInput.Title, "Should have the same title")
		require.Equal(t, len(validInput.Attachments), len(storedFeedback.Attachments), "Should have stored all attachments/fetched all attachments")
		for i := 0; i < len(validInput.Attachments); i++ {
			storedAttch := storedFeedback.Attachments[i]
			insertedAttch := validInput.Attachments[i]
			require.Equal(t, insertedAttch.Name, storedAttch.Name, "Stored Attachment should have the same Name")
			require.Equal(t, insertedAttch.Path, storedAttch.Path, "Stored Attachment should have the same Path")
		}
	})

}

func TestGetFeedbacksWithPagination(t *testing.T) {
	t.Run("Doesn't use FeedbackRepository and 0-length slice and error when pagination limit is 0", func(t *testing.T) {
		t.Cleanup(testing_utils.TruncateTestTables(t, "feedbacks", "attachments"))
		// arrange
		usecase, _, _ := getVars()

		// act
		pagination := feedback_usecases.Pagination{Limit: 0}
		feedbacks, err := usecase.GetFeedbacksWithPagination(context.Background(), pagination)

		// assert
		assert.Equal(t, ZeroLimitPaginationErr, err, "Should return non-nil error when limit is 0")
		assert.Equal(t, 0, len(feedbacks), "Should return nil(0-length) slice when pagination limit is 0")
	})
}
