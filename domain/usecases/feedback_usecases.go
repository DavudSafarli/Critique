package usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type Pagination struct {
	Skip  uint
	Limit uint
}

type FeedbackUsecases interface {
	CreateFeedback(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error)
	GetFeedbackDetails(ctx context.Context, id uint) (models.Feedback, error)
	GetFeedbacksWithPagination(ctx context.Context, pagination Pagination) ([]models.Feedback, error)
}
