package impl

import (
	"context"
	"errors"
	"github.com/DavudSafarli/Critique/domain/contracts"

	"github.com/DavudSafarli/Critique/domain/models"
)

// TagUsecasesImpl is a struct that implements all Tag Usecases
type TagUsecasesImpl struct {
	storage contracts.Storage
}

type ti = TagUsecasesImpl

func (t ti) GetTags(ctx context.Context) ([]models.Tag, error) {
	return t.storage.GetTags(ctx)
}

var errEmptySlice error = errors.New("Passed TagIds is empty")

func (t ti) DeleteTags(ctx context.Context, tagIds []uint) error {
	if len(tagIds) == 0 {
		return errEmptySlice
	}
	return t.storage.RemoveManyTags(ctx, tagIds)
}

func (t ti) CreateTags(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	return tags, t.storage.CreateManyTags(ctx, tags)
}
