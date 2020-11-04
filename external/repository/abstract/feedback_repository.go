package abstract

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

// FeedbackRepository is an interfce
type FeedbackRepository interface {
	// GetPaginated returns records with pagination
	GetPaginated(ctx context.Context, skip uint, limit uint) ([]models.Feedback, error)
	// Find finds and retrieves a single record with the given ID
	Find(ctx context.Context, id uint) (f models.Feedback, err error)
	// Create persists a new Feedback to the database and returns newly inserted Feedback
	Create(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error)
	// UpdateTagIDs just panics right now, but will update "tag_id"s of feedbacks from x to y.
	UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error
}
