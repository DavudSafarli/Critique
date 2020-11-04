package impl

import (
	"context"

	"github.com/DavudSafarli/Critique/external/repository/abstract"

	"github.com/DavudSafarli/Critique/domain/models"
)

// TagUsecasesImpl is a struct that implements all Tag Usecases
type TagUsecasesImpl struct {
	tagRepository abstract.TagRepository
}

type ti = TagUsecasesImpl

func (t ti) GetTags(ctx context.Context) ([]models.Tag, error) {
	return t.tagRepository.Get(ctx)
}

func (t ti) DeleteTags(ctx context.Context, tagIds []uint) error {
	return t.tagRepository.RemoveMany(ctx, tagIds)
}

func (t ti) CreateTags(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	return t.tagRepository.CreateMany(ctx, tags)
}
