package impl

import (
	"context"

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

// CreateAttachment does what it says
func (a AttachmentUsecasesImpl) CreateAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) (attchs []models.Attachment, err error) {
	return a.attachmentRepository.CreateMany(ctx, attachments, feedbackID)
}

func (a AttachmentUsecasesImpl) GetAttachments(ctx context.Context, feedbackID uint) (attchs []models.Attachment, err error) {
	return a.attachmentRepository.GetByFeedbackID(ctx, feedbackID)
}
