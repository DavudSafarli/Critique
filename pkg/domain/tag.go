package domain

import "context"

// Tag is a domain model
type Tag struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TagRepository is the contract that all implementations must implement
type TagRepository interface {
	CreateMany(ctx context.Context, tags []Tag) ([]Tag, error)
	Get(ctx context.Context) ([]Tag, error)
	RemoveMany(ctx context.Context, tagIDs []uint) error
}
