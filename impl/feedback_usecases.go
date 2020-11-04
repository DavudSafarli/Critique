package impl

import (
	"context"

	"github.com/DavudSafarli/Critique/external/repository/abstract"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
)

// FeedbackValidator an an interface for validating Feedback model
type FeedbackValidator interface {
	Validate(feedback models.Feedback) error
}

type fi = FeedbackUsecasesImpl

// FeedbackUsecasesImpl is a struct that implements all Feedback Usecases
type FeedbackUsecasesImpl struct {
	feedbackRepository abstract.FeedbackRepository
	validator          FeedbackValidator
}

// NewFeedbackUsecasesImpl creates new FeedbackUsecasesImpl
func NewFeedbackUsecasesImpl(repo abstract.FeedbackRepository, validator FeedbackValidator) feedback_usecases.FeedbackUsecases {
	return FeedbackUsecasesImpl{
		feedbackRepository: repo,
		validator:          validator,
	}
}

func (g fi) CreateFeedback(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error) {
	if err := g.validator.Validate(feedback); err != nil {
		return f, err
	}
	return g.feedbackRepository.Create(ctx, feedback)
}

func (g fi) GetFeedbackDetails(ctx context.Context, id uint) (models.Feedback, error) {
	return g.feedbackRepository.Find(ctx, id)
}

func (g fi) GetFeedbacksWithPagination(ctx context.Context, pagination feedback_usecases.Pagination) ([]models.Feedback, error) {
	return g.feedbackRepository.GetPaginated(ctx, pagination.Skip, pagination.Limit)
}
