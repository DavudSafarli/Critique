package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTX struct {
	mock.Mock
}

func (m MockTX) BeginTx(ctx context.Context) (context.Context, error) {
	args := m.Called(ctx)
	return args.Get(0).(context.Context), args.Error(1)
}

func (m MockTX) CommitTx(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m MockTX) RollbackTx(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
