package storage

import "context"

// Store abstracts object storage operations.
type Store interface {
	// GenerateUploadURL returns a presigned upload URL and the future public URL for the given filename.
	GenerateUploadURL(ctx context.Context, filename string) (uploadURL string, publicURL string, err error)
	// PublicURL returns the public URL for a stored key.
	PublicURL(ctx context.Context, key string) string
}
