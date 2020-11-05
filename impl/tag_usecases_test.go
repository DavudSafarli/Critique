package impl

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/external/repository/mocks"
	"github.com/bmizerany/assert"
)

func getTagMocksAndUsecase(t *testing.T) (TagUsecasesImpl, *mocks.MockTagRepository) {
	t.Helper()
	mockRepo := new(mocks.MockTagRepository)
	usecase := TagUsecasesImpl{mockRepo}
	return usecase, mockRepo
}

func TestDeleteTags(t *testing.T) {
	t.Run("Doesn't use TagRepository and returns error when slice is empty", func(t *testing.T) {
		t.Parallel()
		// arrange
		usecase, _ := getTagMocksAndUsecase(t)

		// act
		var tagIds []uint = nil
		err := usecase.DeleteTags(context.Background(), tagIds)

		// assert
		assert.NotEqual(t, nil, err, "Should return non-nil error when slice is empty")
	})
}
