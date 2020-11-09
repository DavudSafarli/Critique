package abstract

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

// AttachmentRepository is an interfce
type AttachmentRepository interface {
	OnePhaseCommitProtocol
	// CreateMany persists new Tags into the database
	CreateMany(ctx context.Context, attachments []models.Attachment) ([]models.Attachment, error)
}
