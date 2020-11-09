package attachment_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type AttachmentCreator interface {
	CreateAttachment(ctx context.Context, attachments []models.Attachment) ([]models.Attachment, error)
}
