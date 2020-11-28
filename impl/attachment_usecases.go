package impl

import (
	"context"
	"errors"

	"github.com/DavudSafarli/Critique/domain/usecases/attachment_usecases"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

// AttachmentUsecasesImpl is a struct that implements all Attachment Usecases
type AttachmentUsecasesImpl struct {
	attachmentRepository abstract.AttachmentRepository
}
type ai = AttachmentUsecasesImpl

// NewAttachmentUsecases creates a new AttachmentUsecases
func NewAttachmentUsecases(attachmentRepo abstract.AttachmentRepository) attachment_usecases.AttachmentUsecases {
	return AttachmentUsecasesImpl{
		attachmentRepository: attachmentRepo,
	}
}

var CreateFeedbackErr = errors.New("create-attachment-err")

// CreateAttachment does what it says
func (a AttachmentUsecasesImpl) CreateAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) (attchs []models.Attachment, err error) {
	if len(attachments) == 0 {
		return nil, CreateFeedbackErr
	}
	return attachments, a.attachmentRepository.CreateMany(ctx, attachments, feedbackID)
}
