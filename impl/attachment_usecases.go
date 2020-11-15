package impl

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository/abstract"
)

// AttachmentUsecasesImpl is a struct that implements all Attachment Usecases
type AttachmentUsecasesImpl struct {
	attachmentRepository abstract.AttachmentRepository
}

func (a AttachmentUsecasesImpl) CreateAttachment(ctx context.Context, attachments []models.Attachment) ([]models.Attachment, error) {
	panic("implement me")
}

type ai = AttachmentUsecasesImpl
