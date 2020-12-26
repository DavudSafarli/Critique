package impl

import (
	"context"
	"errors"

	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/domain/models"
)

// AttachmentUsecasesImpl is a struct that implements all Attachment Usecases
type AttachmentUsecasesImpl struct {
	storage contracts.Storage
}
type ai = AttachmentUsecasesImpl

// NewAttachmentUsecases creates a new AttachmentUsecases
func NewAttachmentUsecases(storage contracts.Storage) AttachmentUsecasesImpl {
	return AttachmentUsecasesImpl{
		storage: storage,
	}
}

var CreateFeedbackErr = errors.New("create-attachment-err")

// CreateAttachment does what it says
func (a AttachmentUsecasesImpl) CreateAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) (attchs []models.Attachment, err error) {
	if len(attachments) == 0 {
		return nil, CreateFeedbackErr
	}
	return attachments, a.storage.CreateManyAttachments(ctx, attachments, feedbackID)
}
