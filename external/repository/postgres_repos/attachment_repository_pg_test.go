package postgres_repos

import (
	"context"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

func TestAttachmentRepository(t *testing.T) {
	storage := vars.storage
	attchrepo := AttachmentRepository{storage}
	fdbkRepo := FeedbackRepository{storage}

	abstract.TestAttachmentRepositoryBehaviour(t, attchrepo, fdbkRepo, getCleanupFuncForAttachments(storage), AttchFuncsForBehaviourTest{storage})
}

// AttchFuncsForBehaviourTest is a struct that implements AttchRequiredFuncs interface, so that contract test suite can do its work
type AttchFuncsForBehaviourTest struct {
	storage *Storage
}

func (f AttchFuncsForBehaviourTest) GetAllAttachments() ([]models.Attachment, error) {
	rows, err := f.storage.DB.Query(context.Background(), "SELECT id, name, path, feedback_id from attachments")
	if err != nil {
		panic(err)
	}
	got := []models.Attachment{}
	for rows.Next() {
		var r models.Attachment
		err = rows.Scan(&r.ID, &r.Name, &r.Path, &r.FeedbackID)
		if err != nil {
			panic(err)
		}
		got = append(got, r)
	}
	return got, nil
}

func getCleanupFuncForAttachments(storage *Storage) func() error {
	return func() error {
		_, err := storage.DB.Exec(context.Background(), "DELETE from attachments")
		if err != nil {
			return err
		}
		_, err = storage.DB.Exec(context.Background(), "DELETE from feedbacks")
		if err != nil {
			return err
		}
		return nil
	}
}
