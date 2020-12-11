package contracts

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type Storage interface {
	AttachmentRepository
	FeedbackRepository
	TagRepository
	OnePhaseCommitProtocol
}

type AttachmentRepository interface {
	//OnePhaseCommitProtocol
	// CreateMany persists new Tags into the database
	CreateManyAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) error
	GetAttachmentsByFeedbackID(ctx context.Context, feedbackID uint) ([]models.Attachment, error)
}

// FeedbackRepository is an interfce
type FeedbackRepository interface {
	//OnePhaseCommitProtocol
	// GetPaginated returns records with pagination
	GetFeedbacksPaginated(ctx context.Context, skip uint, limit uint) ([]models.Feedback, error)
	// Find finds and retrieves a single record with the given ID
	FindFeedback(ctx context.Context, id uint) (f models.Feedback, err error)
	// Create persists a new Feedback to the database and returns newly inserted Feedback
	CreateFeedback(ctx context.Context, feedback *models.Feedback) (err error)
	// UpdateTagIDs just panics right now, but will update "tag_id"s of feedbacks from x to y.
	UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error
}

// TagRepository is an interfce
type TagRepository interface {
	//OnePhaseCommitProtocol
	// CreateMany persists new Tags into the database
	CreateManyTags(ctx context.Context, tags []models.Tag) error
	// Get returns all Tags
	GetTags(ctx context.Context) ([]models.Tag, error)
	// RemoveMany removes Tags of given tagIDs from database
	RemoveManyTags(ctx context.Context, tagIDs []uint) error
}

type OnePhaseCommitProtocol interface {
	// BeginTx creates a context with a transaction.
	// All statements that receive this context should be executed within the given transaction in the context.
	// After a BeginTx command will be executed in a single transaction until an explicit COMMIT or ROLLBACK is given.
	BeginTx(ctx context.Context) (context.Context, error)
	// Commit commits the current transaction.
	// All changes made by the transaction become visible to others and are guaranteed to be durable if a crash occurs.
	CommitTx(context.Context) error
	// Rollback rolls back the current transaction and causes all the updates made by the transaction to be discarded.
	RollbackTx(context.Context) error
}
