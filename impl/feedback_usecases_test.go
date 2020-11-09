package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
	"github.com/DavudSafarli/Critique/external/repository/mocks"
	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/mock"
)

type MockFeedbackValidator struct {
	mock.Mock
}

func (m *MockFeedbackValidator) Validate(feedback models.Feedback) error {
	args := m.Called(feedback)
	return args.Error(0)
}

func getFeedbackMocksAndUsecase(t *testing.T) (FeedbackUsecasesImpl, *mocks.MockFeedbackRepository, *MockFeedbackValidator) {
	t.Helper()
	mockRepo := new(mocks.MockFeedbackRepository)
	mockAttchRepo := new(mocks.MockAttachmentRepository)
	mockValidator := new(MockFeedbackValidator)
	usecase := FeedbackUsecasesImpl{mockRepo, mockAttchRepo, mockValidator}
	return usecase, mockRepo, mockValidator
}
func TestCreateFeedback(t *testing.T) {
	t.Run("Calls FeedbackRepo.Create on successful case", func(t *testing.T) {
		t.Parallel()
		// arrange
		usecase, mockRepo, mockValidator := getFeedbackMocksAndUsecase(t)
		mockValidator.On("Validate", mock.Anything).Return(nil)
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(models.Feedback{}, nil)

		// act
		f := models.Feedback{Title: "Salam"}
		usecase.CreateFeedback(context.Background(), f)

		// assert
		mockRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Calls Validate and returns error on Validation fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		returnedErr := errors.New("salam")
		usecase, _, mockValidator := getFeedbackMocksAndUsecase(t)
		mockValidator.On("Validate", mock.Anything).Return(returnedErr)

		// act
		f := models.Feedback{}
		_, err := usecase.CreateFeedback(context.Background(), f)

		// assert
		mockValidator.AssertCalled(t, "Validate", f)
		assert.Equal(t, returnedErr, err)
		assert.NotEqual(t, nil, err, "Should return error if validate fails")
	})

	t.Run("Calls repository on successful validation and returns repository results", func(t *testing.T) {
		t.Parallel()
		// arrange
		argumentPassed := models.Feedback{}
		usecase, mockRepo, mockValidator := getFeedbackMocksAndUsecase(t)
		mockValidator.On("Validate", argumentPassed).Return(nil)
		mockRepo.On("Create", mock.Anything, argumentPassed).Return(argumentPassed, errors.New("mock"))

		// act
		f, err := usecase.CreateFeedback(context.Background(), argumentPassed)

		// assert
		mockRepo.AssertCalled(t, "Create", mock.Anything, argumentPassed)
		assert.Equal(t, f, argumentPassed)
		assert.Equal(t, errors.New("mock"), err, "Should return error if validate fails")
	})
}

func TestGetFeedbacksWithPagination(t *testing.T) {
	t.Run("Doesn't use FeedbackRepository and 0-length slice and error when pagination limit is 0", func(t *testing.T) {
		// arrange
		usecase, _, _ := getFeedbackMocksAndUsecase(t)

		// act
		pagination := feedback_usecases.Pagination{Limit: 0}
		feedbacks, err := usecase.GetFeedbacksWithPagination(context.Background(), pagination)

		// assert
		assert.Equal(t, ZeroLimitPaginationErr, err, "Should return non-nil error when limit is 0")
		assert.Equal(t, 0, len(feedbacks), "Should return nil(0-length) slice when pagination limit is 0")
	})
}
