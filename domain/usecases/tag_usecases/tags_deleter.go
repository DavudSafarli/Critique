package feedback_usecases

import (
	"context"
)

type TagsDeleter interface {
	DeleteTags(ctx context.Context, tagIds []uint) error
}
