package abstract

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bmizerany/assert"

	"github.com/DavudSafarli/Critique/domain/models"

	"github.com/DavudSafarli/Critique/util"
)

type AttchRequiredFuncs interface {
	GetAllAttachments() ([]models.Attachment, error)
}

// TestAttachmentRepositoryBehaviour does what its name says
func TestAttachmentRepositoryBehaviour(t *testing.T, abstractRepo AttachmentRepository,
	abstractFeedbackRepo FeedbackRepository, cleanupFunc func() error, funcs AttchRequiredFuncs) {
	t.Run("Create attachments successfully", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
		testCreateMany(t, abstractRepo, abstractFeedbackRepo, funcs)
	})
}

// TestCreate does what its name says
func testCreateMany(t *testing.T, attchRepo AttachmentRepository, fdbkRepo FeedbackRepository, funcs AttchRequiredFuncs) {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	f, _ := fdbkRepo.Create(ctx, models.Feedback{
		Title:     fmt.Sprint("Title"),
		Body:      fmt.Sprint("Body"),
		CreatedBy: "uniqueUserIdentifier",
		CreatedAt: uint(time.Now().Unix()),
	})

	// act
	attchs := []models.Attachment{}
	var numOfInsertedRows int = 10
	for i := 0; i < numOfInsertedRows; i++ {
		attchs = append(attchs, models.Attachment{
			Name:       fmt.Sprint("Name", i),
			Path:       fmt.Sprint("Name", i),
			FeedbackID: f.ID,
		})
	}
	returnedAttchs, err := attchRepo.CreateMany(ctx, attchs)
	if err != nil {
		t.Fatalf("Failed to Create-Many Attachments : %s", err)
	}

	// assert
	attchsInDb, _ := funcs.GetAllAttachments()
	assert.Equal(t, numOfInsertedRows, len(attchsInDb), "Database should contain the same number of rows")

	assert.Equal(t, numOfInsertedRows, len(returnedAttchs), "CreateMany should return the same num. of attachments")
}
