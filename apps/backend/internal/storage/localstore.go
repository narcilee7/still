package storage

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/google/uuid"
)

// LocalStore stores uploads on the local filesystem and serves them via HTTP.
type LocalStore struct {
	baseURL   string
	uploadDir string
}

// NewLocalStore creates a new LocalStore.
func NewLocalStore(baseURL, uploadDir string) *LocalStore {
	return &LocalStore{
		baseURL:   baseURL,
		uploadDir: uploadDir,
	}
}

// GenerateUploadURL returns a URL that the local upload handler accepts.
func (s *LocalStore) GenerateUploadURL(ctx context.Context, filename string) (string, string, error) {
	key := uuid.New().String() + path.Ext(filename)
	uploadURL, err := url.JoinPath(s.baseURL, "_upload", key)
	if err != nil {
		return "", "", fmt.Errorf("build upload url failed: %w", err)
	}
	publicURL, err := url.JoinPath(s.baseURL, "_files", key)
	if err != nil {
		return "", "", fmt.Errorf("build public url failed: %w", err)
	}
	return uploadURL, publicURL, nil
}

// PublicURL returns the public URL for a stored key.
func (s *LocalStore) PublicURL(ctx context.Context, key string) string {
	u, _ := url.JoinPath(s.baseURL, "_files", key)
	return u
}

// UploadDir returns the directory where uploaded files are stored.
func (s *LocalStore) UploadDir() string {
	return s.uploadDir
}
