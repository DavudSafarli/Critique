package usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type AttachmentUsecases interface {
	CreateAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) ([]models.Attachment, error)
}
