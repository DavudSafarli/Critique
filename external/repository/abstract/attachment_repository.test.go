package abstract

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	testGetByFeedbackID(t, attchRepo, abstractFeedbackRepo, cleanupFunc)
}

// testCreateManyAndGetAll does what its name says
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

// testGetByFeedbackID makes BDD style testing
func testGetByFeedbackID(t *testing.T, attchRepo AttachmentRepositoryTester, fdbkRepo FeedbackRepository, cleanupFunc func()) {
	spec := testcase.NewSpec(t)
	spec.After(func(t *testcase.T) { cleanupFunc() })

	spec.Describe(`#GetByFeedbackID`, func(s *testcase.Spec) {
		feedbackID := testcase.Var{Name: `feedbackID`}
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			return attchRepo.GetByFeedbackID(context.Background(), feedbackID.Get(t).(uint))
		}
		s.When(`there are attachments`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				f, err := fdbkRepo.Create(context.Background(), models.Feedback{Title: "T", Body: "T", CreatedBy: "T", CreatedAt: uint(time.Now().Unix())})
				require.Nil(t, err)
				feedbackID.Set(t, f.ID) // <---------
				_, err = attchRepo.CreateMany(context.Background(), []models.Attachment{{Name: "A1", Path: "P1"}, {Name: "A2", Path: "P2"}}, f.ID)
				require.Nil(t, err)

				feedback2 := models.Feedback{Title: "T2", Body: "T2", CreatedBy: "T2", CreatedAt: uint(time.Now().Unix())}
				_, err = attchRepo.CreateMany(context.Background(), []models.Attachment{{Name: "Z", Path: "Z"}}, feedback2.ID)
				require.Nil(t, err)
			})
			s.Then(`should find them`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err, "No error pls")
				require.Len(t, attchs, 2)
			})
		})
	})

}
