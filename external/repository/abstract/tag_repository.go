package abstract

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

// TagRepository is an interfce
type TagRepository interface {
	OnePhaseCommitProtocol
	// CreateMany persists new Tags into the database
	CreateMany(ctx context.Context, tags []models.Tag) ([]models.Tag, error)
	// Get returns all Tags
	Get(ctx context.Context) ([]models.Tag, error)
	// RemoveMany removes Tags of given tagIDs from database
	RemoveMany(ctx context.Context, tagIDs []uint) error
}
