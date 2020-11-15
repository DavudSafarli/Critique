package abstract

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
)

// TestFeedbackRepositoryBehaviour does what its name says
func TestFeedbackRepositoryBehaviour(t *testing.T, abstractRepo FeedbackRepository, cleanupFunc func()) {
	t.Run("Creates Feedbacks and Get them successfully", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		testCreatesAndGetAll(t, abstractRepo)
	})
	t.Run("GetPaginated returns in correct size and order", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		testGetPaginated(t, abstractRepo)
	})
	t.Run("can Find Created Feedbacks", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		testFind(t, abstractRepo)
	})
	t.Run("test UpdateTagIDs", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		testUpdateTagIDs(t, abstractRepo)
	})
}

// testCreate does what its name says
func testCreatesAndGetAll(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 10
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	var skip uint = 0
	var limit uint = uint(numOfInsertedRows * 2)
	selectedFeedbacks, err := abstractRepo.GetPaginated(ctx, skip, limit)
	require.Nil(t, err)
	isEqual := reflect.DeepEqual(insertedFeedbacks, selectedFeedbacks)
	require.True(t, isEqual, "Inserted Feedbacks are not the same as selected ones database")
}

// testGetPaginated does what its name says
func testGetPaginated(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 10
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)
	require.Equal(t, len(insertedFeedbacks), numOfInsertedRows, len(insertedFeedbacks), "is not equal", numOfInsertedRows)
	var skip int = 0
	var limit uint = 3
	for skip < numOfInsertedRows {
		selectedFeedbacks, err := abstractRepo.GetPaginated(ctx, uint(skip), limit)
		require.Nil(t, err)
		lengthShouldBe := int(math.Min(float64(numOfInsertedRows-skip), float64(limit)))
		require.Equal(t, lengthShouldBe, len(selectedFeedbacks), "GetPaginated did not return correct size slice")

		// assert if selected rows are the same as the inserted rows
		isEqual := reflect.DeepEqual(insertedFeedbacks[skip:skip+lengthShouldBe], selectedFeedbacks)
		require.True(t, isEqual, "Inserted Feedbacks are not the same as selected ones database. ")
		skip += 3
	}

}

// testFind does what its name says
func testFind(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 4
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	for _, f := range insertedFeedbacks {
		foundFeedback, err := abstractRepo.Find(ctx, f.ID)
		require.Nil(t, err, "Failed to perform Find")

		require.Equal(t, f.ID, foundFeedback.ID, "Failed to perform Find")
	}
}

// testUpdateTagIDs does what its name says
func testUpdateTagIDs(t *testing.T, abstractRepo FeedbackRepository) {
	// This is not required right now
}

func insertNFeedbacks(ctx context.Context, t *testing.T, n int, abstractRepo FeedbackRepository) []models.Feedback {
	insertedFeedbacks := make([]models.Feedback, 0, n)
	for i := 0; i < n; i++ {
		feedback, err := abstractRepo.Create(ctx, models.Feedback{
			Title:     fmt.Sprint("Title", i),
			Body:      fmt.Sprint("Body", i),
			CreatedBy: "uniqueUserIdentifier",
			CreatedAt: uint(time.Now().Unix()),
		})
		require.Nil(t, err, "Create: Failed to create Feedback")
		insertedFeedbacks = append(insertedFeedbacks, feedback)
	}
	//time.Sleep(5 * time.Second)
	return insertedFeedbacks
}
