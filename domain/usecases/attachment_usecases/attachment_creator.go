package attachment_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type AttachmentCreator interface {
	CreateAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) ([]models.Attachment, error)
}
