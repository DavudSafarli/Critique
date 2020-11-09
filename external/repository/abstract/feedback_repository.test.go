package abstract

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/util"
)

// TestFeedbackRepositoryBehaviour does what its name says
func TestFeedbackRepositoryBehaviour(t *testing.T, abstractRepo FeedbackRepository, cleanupFunc func() error) {
	t.Run("Test Create", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
		TestCreate(t, abstractRepo)
	})
	t.Run("Test GetPaginated", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
		TestGetPaginated(t, abstractRepo)
	})
	t.Run("Test Find", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
		TestFind(t, abstractRepo)
	})
	t.Run("Test UpdateTagIDs", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
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
	selectedFeedbacks, _ := abstractRepo.GetPaginated(ctx, skip, limit)

	if reflect.DeepEqual(insertedFeedbacks, selectedFeedbacks) == false {
		t.Fatal("Inserted Feedbacks are not the same as selected ones database")
	}

}

// TestGetPaginated does what its name says
func TestGetPaginated(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 10
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	var skip int = 0
	var limit uint = 3
	for skip < numOfInsertedRows {
		selectedFeedbacks, _ := abstractRepo.GetPaginated(ctx, uint(skip), limit)
		lengthShouldBe := int(math.Min(float64(numOfInsertedRows-skip), float64(limit)))
		if len(selectedFeedbacks) != lengthShouldBe {
			t.Fatal("GetPaginated did not return correct size slice")
		}

		// assert if selected rows are the same as the inserted rows
		if reflect.DeepEqual(insertedFeedbacks[skip:skip+lengthShouldBe], selectedFeedbacks) == false {
			t.Fatal("Inserted Feedbacks are not the same as selected ones database")
		}
		skip += 3
	}

}

// TestFind does what its name says
func TestFind(t *testing.T, abstractRepo FeedbackRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	var numOfInsertedRows int = 4
	insertedFeedbacks := insertNFeedbacks(ctx, t, numOfInsertedRows, abstractRepo)

	for _, f := range insertedFeedbacks {
		foundFeedback, err := abstractRepo.Find(ctx, f.ID)
		if err != nil {
			t.Fatalf("Failed to perform Find: %s", err)
		}
		if f.ID != foundFeedback.ID {
			t.Fatalf("Retrieved Feedback is not the same as the inserted one")
		}
	}
}

// TestUpdateTagIDs does what its name says
func TestUpdateTagIDs(t *testing.T, abstractRepo FeedbackRepository) {
	// This is not required right now
}

func insertNFeedbacks(ctx context.Context, t *testing.T, n int, abstractRepo FeedbackRepository) []models.Feedback {
	t.Helper()
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
	return insertedFeedbacks
}
