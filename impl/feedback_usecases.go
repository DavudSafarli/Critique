package impl

import (
	"context"
	"errors"

	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/domain/usecases"

	"github.com/DavudSafarli/Critique/domain/models"
)

type fi = FeedbackUsecasesImpl

// FeedbackUsecasesImpl is a struct that implements all Feedback Usecases
type FeedbackUsecasesImpl struct {
	storage contracts.Storage
}

// NewFeedbackUsecasesImpl creates new FeedbackUsecasesImpl
func NewFeedbackUsecasesImpl(storage contracts.Storage) usecases.FeedbackUsecases {
	return &FeedbackUsecasesImpl{
		storage: storage,
	}
}

var createFeedbackErr = errors.New("create-feedback-err")

// TODO: add integration test for validating Atomicity of CreateFeedback. If CreateAttchmnt fails, then feedback should not persist either
func (g *fi) CreateFeedback(ctx context.Context, feedback models.Feedback) (emptyFeedback models.Feedback, err error) {
	if err = feedback.Validate(); err != nil {
		return
	}
	ctx, err = g.storage.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() { g.commitOrRollback(ctx, err) }()

	if err = g.storage.CreateFeedback(ctx, &feedback); err != nil {
		return emptyFeedback, createFeedbackErr
	}
	if feedback.Attachments == nil {
		return feedback, nil
	}
	err = g.storage.CreateManyAttachments(ctx, feedback.Attachments, feedback.ID)
	if err != nil {
		return
	}
	return feedback, nil
}

func (g *fi) GetFeedbackDetails(ctx context.Context, id uint) (models.Feedback, error) {
	return g.storage.FindFeedback(ctx, id)
}

func (g *fi) GetFeedbacksWithPagination(ctx context.Context, pagination usecases.Pagination) ([]models.Feedback, error) {
	if pagination.Limit == 0 {
		return nil, ZeroLimitPaginationErr
	}

	return g.storage.GetFeedbacksPaginated(ctx, pagination.Skip, pagination.Limit)
}

func (g *fi) commitOrRollback(ctx context.Context, err error) {
	if err != nil {
		err = g.storage.RollbackTx(ctx)
		return
	}
	err = g.storage.CommitTx(ctx)
}
