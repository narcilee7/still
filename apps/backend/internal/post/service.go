package post

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// Service implements PostService.
type Service struct{}

// NewService creates a new post service.
func NewService() *Service {
	return &Service{}
}

// CreatePost creates a new post.
func (s *Service) CreatePost(ctx context.Context, req *connect.Request[stillv1.CreatePostRequest]) (*connect.Response[stillv1.CreatePostResponse], error) {
	return connect.NewResponse(&stillv1.CreatePostResponse{}), nil
}

// GetPost returns a post by ID.
func (s *Service) GetPost(ctx context.Context, req *connect.Request[stillv1.GetPostRequest]) (*connect.Response[stillv1.GetPostResponse], error) {
	return connect.NewResponse(&stillv1.GetPostResponse{}), nil
}
