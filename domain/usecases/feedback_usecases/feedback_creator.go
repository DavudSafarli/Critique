package feedback_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

// bla bla
type FeedbackCreator interface {
	// Validates Feedback
	//
	// Creates Feedback and determines ID
	//
	// Delivers the newly created Feedback
	CreateFeedback(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error)
}
