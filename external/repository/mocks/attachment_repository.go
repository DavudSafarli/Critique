package mocks

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockAttachmentRepository is a mock
type MockAttachmentRepository struct {
	mock.Mock
	*MockTX
}

func (m MockAttachmentRepository) CreateMany(ctx context.Context, attachments []models.Attachment) ([]models.Attachment, error) {
	args := m.Called(ctx, attachments)
	return args.Get(0).([]models.Attachment), args.Error(1)
}
