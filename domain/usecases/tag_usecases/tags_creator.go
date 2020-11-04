package feedback_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type TagsCreator interface {
	CreateTags(ctx context.Context, tags []models.Tag) ([]models.Tag, error)
}
