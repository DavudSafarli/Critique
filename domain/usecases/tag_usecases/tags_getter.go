package feedback_usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type TagsGetter interface {
	GetTags(ctx context.Context) ([]models.Tag, error)
}
