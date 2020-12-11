package models

import (
	"strings"
)

// Feedback is a models model
type Feedback struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	CreatedAt uint   `json:"created_at,omitempty"`

	Attachments []Attachment `json:"attachments,omitempty"`
}

type StandardError string

func (e StandardError) Error() string {
	return string(e)
}

const (
	INVALID_FEEDBACK StandardError = "InvalidFeedback"
)

func (f Feedback) Validate() error {
	if len(strings.TrimSpace(f.Title)) == 0 {
		return INVALID_FEEDBACK
	}
	return nil
}

// FeedbackRepository is the contract that all implementations must implement
//type FeedbackRepository interface {
//	GetFeedbacksPaginated(ctx context.Context, skip uint, limit uint) ([]Feedback, error)
//	Find(ctx context.Context, id uint) (Feedback, error)
//	Create(ctx context.Context, feedback Feedback) (Feedback, error)
//	UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error
//}
