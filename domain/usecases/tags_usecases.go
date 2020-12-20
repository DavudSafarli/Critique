package usecases

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
)

type TagsUsecases interface {
	CreateTags(ctx context.Context, tags []models.Tag) ([]models.Tag, error)
	GetTags(ctx context.Context) ([]models.Tag, error)
	DeleteTags(ctx context.Context, tagIds []uint) error
}
