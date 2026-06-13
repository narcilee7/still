package storage

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// Service implements StorageService.
type Service struct {
	store Store
}

// NewService creates a new storage service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// GetUploadURL returns a presigned upload URL and the future public URL.
func (s *Service) GetUploadURL(ctx context.Context, req *connect.Request[stillv1.GetUploadURLRequest]) (*connect.Response[stillv1.GetUploadURLResponse], error) {
	if strings.TrimSpace(req.Msg.Filename) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("filename is required"))
	}
	if strings.TrimSpace(req.Msg.ContentType) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("content_type is required"))
	}
	uploadURL, publicURL, err := s.store.GenerateUploadURL(ctx, req.Msg.Filename, req.Msg.ContentType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&stillv1.GetUploadURLResponse{
		UploadUrl: uploadURL,
		PublicUrl: publicURL,
	}), nil
}
