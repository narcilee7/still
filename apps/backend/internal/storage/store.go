package storage

import "context"

// Store abstracts object storage operations.
type Store interface {
	GenerateUploadURL(ctx context.Context, key string) (string, error)
	PublicURL(ctx context.Context, key string) string
}

// StubStore is a placeholder implementation.
type StubStore struct{}

func (s *StubStore) GenerateUploadURL(ctx context.Context, key string) (string, error) {
	return "https://example.com/upload/" + key, nil
}

func (s *StubStore) PublicURL(ctx context.Context, key string) string {
	return "https://example.com/" + key
}
