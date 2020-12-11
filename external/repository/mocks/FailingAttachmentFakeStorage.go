package mocks

import (
	"context"
	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockAttachmentRepository is a mock
type FailingAttachmentFakeStorage struct {
	mock.Mock
	contracts.Storage
}

func (m FailingAttachmentFakeStorage) CreateManyAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) error {
	args := m.Called(ctx, attachments, feedbackID)
	return args.Error(0)
}

//func (m MockStorage) GetAttachmentsByFeedbackID(ctx context.Context, feedbackID uint) ([]models.Attachment, error) {
//	panic("implement me")
//}
