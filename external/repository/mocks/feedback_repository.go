package mocks

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockFeedbackRepository is a mock
type MockFeedbackRepository struct {
	mock.Mock
}

// GetPaginated is a mock
func (m *MockFeedbackRepository) GetPaginated(ctx context.Context, skip uint, limit uint) ([]models.Feedback, error) {
	args := m.Called(ctx, skip, limit)
	return args.Get(0).([]models.Feedback), args.Error(1)
}

// Find is a mock
func (m *MockFeedbackRepository) Find(ctx context.Context, id uint) (f models.Feedback, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Feedback), args.Error(1)
}

// Create is a mock
func (m *MockFeedbackRepository) Create(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error) {
	args := m.Called(ctx, feedback)
	return args.Get(0).(models.Feedback), args.Error(1)
}

// UpdateTagIDs is a mock
func (m *MockFeedbackRepository) UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error {
	args := m.Called(ctx, tagIDFrom, tagIDTo)
	return args.Error(0)
}
