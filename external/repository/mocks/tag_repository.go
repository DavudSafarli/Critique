package mocks

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/stretchr/testify/mock"
)

// TagRepository is an interfce
type MockTagRepository struct {
	mock.Mock
	*MockTX
}

func (m *MockTagRepository) CreateMany(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	args := m.Called(ctx, tags)
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (m *MockTagRepository) Get(ctx context.Context) ([]models.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (m *MockTagRepository) RemoveMany(ctx context.Context, tagIDs []uint) error {
	args := m.Called(ctx, tagIDs)
	return args.Error(1)
}
