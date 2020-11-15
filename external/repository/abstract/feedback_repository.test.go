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
	t.Run("Test Create", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		TestCreate(t, abstractRepo)
	})
	t.Run("Test GetPaginated", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		TestGetPaginated(t, abstractRepo)
	})
	t.Run("Test Find", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		TestFind(t, abstractRepo)
	})
	t.Run("Test UpdateTagIDs", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		TestUpdateTagIDs(t, abstractRepo)
	})
}

// TestCreate does what its name says
func TestCreate(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 10
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	var skip uint = 0
	var limit uint = 20
	selectedFeedbacks, err := abstractRepo.GetPaginated(ctx, skip, limit)
	require.Nil(t, err)
	if reflect.DeepEqual(insertedFeedbacks, selectedFeedbacks) == false {
		t.Fatal("Inserted Feedbacks are not the same as selected ones database. ", "selectedFeedbacks length = ", len(selectedFeedbacks))
	}

}

// TestGetPaginated does what its name says
func TestGetPaginated(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 100
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)
	require.Equal(t, len(insertedFeedbacks), numOfInsertedRows, len(insertedFeedbacks), "is not equal", numOfInsertedRows)
	var skip int = 0
	var limit uint = 3
	for skip < numOfInsertedRows {
		selectedFeedbacks, err := abstractRepo.GetPaginated(ctx, uint(skip), limit)
		require.Nil(t, err)
		lengthShouldBe := int(math.Min(float64(numOfInsertedRows-skip), float64(limit)))
		if len(selectedFeedbacks) != lengthShouldBe {
			t.Fatal("GetPaginated did not return correct size slice")
		}

		// assert if selected rows are the same as the inserted rows
		if reflect.DeepEqual(insertedFeedbacks[skip:skip+lengthShouldBe], selectedFeedbacks) == false {
			t.Fatal("Inserted Feedbacks are not the same as selected ones database", skip, ":", skip+lengthShouldBe)
		}
		skip += 3
	}

}

// TestFind does what its name says
func TestFind(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 40
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	for _, f := range insertedFeedbacks {
		foundFeedback, err := abstractRepo.Find(ctx, f.ID)
		if err != nil {
			t.Fatalf("Failed to perform Find: %s", err)
		}
		if f.ID != foundFeedback.ID {
			t.Fatal("Retrieved Feedback is not the same as the inserted one", f, foundFeedback)
		}
	}
}

// TestUpdateTagIDs does what its name says
func TestUpdateTagIDs(t *testing.T, abstractRepo FeedbackRepository) {
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
		if err != nil {
			t.Fatalf("Create: Failed to create Feedback: %s", err)
		}
		insertedFeedbacks = append(insertedFeedbacks, feedback)
	}
	//time.Sleep(5 * time.Second)
	return insertedFeedbacks
}
