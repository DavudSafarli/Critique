package feedback_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type Pagination struct {
	Skip  uint
	Limit uint
}
type FeedbackPaginator interface {
	GetFeedbacksWithPagination(ctx context.Context, pagination Pagination) ([]models.Feedback, error)
}
