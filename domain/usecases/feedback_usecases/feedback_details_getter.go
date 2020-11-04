package feedback_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type FeedbackDetailsGetter interface {
	GetFeedbackDetails(ctx context.Context, id uint) (models.Feedback, error)
}
