package abstract

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
)

type AttchRequiredFuncs interface {
	GetAllAttachments() ([]models.Attachment, error)
}

type AttachmentRepositoryTester interface {
	AttachmentRepository
	GetAll(ctx context.Context) ([]models.Attachment, error)
}

// TestAttachmentRepositoryBehaviour does what its name says
func TestAttachmentRepositoryBehaviour(t *testing.T, attchRepo AttachmentRepositoryTester, abstractFeedbackRepo FeedbackRepository, cleanupFunc func()) {
	t.Run("Creates attachments and GetAll them successfully", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		testCreateManyAndGetAll(t, attchRepo, abstractFeedbackRepo)
	})
}

// TestCreate does what its name says
func testCreateManyAndGetAll(t *testing.T, attchRepo AttachmentRepositoryTester, fdbkRepo FeedbackRepository) {
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
	returnedAttchs, err := attchRepo.CreateMany(ctx, attchs, f.ID)
	require.Nil(t, err, "Failed to Create-Many Attachments")

	// assert
	attchsInDb, err := attchRepo.GetAll(ctx)
	assert.Nil(t, err, "GetAll Attachments Should not return error")
	assert.Equal(t, numOfInsertedRows, len(attchsInDb), "Database should contain the same number of rows")

	assert.Equal(t, numOfInsertedRows, len(returnedAttchs), "CreateMany should return the same num. of attachments")
}
