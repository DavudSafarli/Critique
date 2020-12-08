package impl

import (
	"context"
	"errors"
	"github.com/DavudSafarli/Critique/external/repository/abstract"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
)

type fi = FeedbackUsecasesImpl

// FeedbackUsecasesImpl is a struct that implements all Feedback Usecases
type FeedbackUsecasesImpl struct {
	feedbackRepository abstract.FeedbackRepository
	attchRepo          abstract.AttachmentRepository
	txer               abstract.OnePhaseCommitProtocol
}

// NewFeedbackUsecasesImpl creates new FeedbackUsecasesImpl
func NewFeedbackUsecasesImpl(repo abstract.FeedbackRepository, attchRepo abstract.AttachmentRepository) feedback_usecases.FeedbackUsecases {
	return FeedbackUsecasesImpl{
		feedbackRepository: repo,
		attchRepo:          attchRepo,
		txer:               repo,
	}
}

var createFeedbackErr = errors.New("create-feedback-err")

// TODO: add integration test for validating Atomicity of CreateFeedback. If CreateAttchmnt fails, then feedback should not persist either
func (g fi) CreateFeedback(ctx context.Context, feedback models.Feedback) (emptyFeedback models.Feedback, err error) {
	if err = feedback.Validate(); err != nil {
		return
	}
	ctx, err = g.txer.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() { g.commitOrRollback(ctx, err) }()

	if err = g.feedbackRepository.Create(ctx, &feedback); err != nil {
		return emptyFeedback, createFeedbackErr
	}
	if feedback.Attachments == nil {
		return feedback, nil
	}
	err = g.attchRepo.CreateMany(ctx, feedback.Attachments, feedback.ID)
	if err != nil {
		return
	}
	return feedback, nil
}

func (g fi) GetFeedbackDetails(ctx context.Context, id uint) (models.Feedback, error) {
	return g.feedbackRepository.Find(ctx, id)
}

func (g fi) GetFeedbacksWithPagination(ctx context.Context, pagination feedback_usecases.Pagination) ([]models.Feedback, error) {
	if pagination.Limit == 0 {
		return nil, ZeroLimitPaginationErr
	}

	return g.feedbackRepository.GetPaginated(ctx, pagination.Skip, pagination.Limit)
}

func (g fi) commitOrRollback(ctx context.Context, err error) {
	if err != nil {
		err = g.txer.RollbackTx(ctx)
		return
	}
	err = g.txer.CommitTx(ctx)
}
