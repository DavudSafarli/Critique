package domain

import "context"

// Feedback is a domain model
type Feedback struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	CreatedAt uint   `json:"created_at,omitempty"`
}

// FeedbackRepository is the contract that all implementations must implement
type FeedbackRepository interface {
	GetPaginated(ctx context.Context, skip uint, limit uint) ([]Feedback, error)
	Find(ctx context.Context, id uint) (Feedback, error)
	Create(ctx context.Context, feedback Feedback) (Feedback, error)
	UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error
}
