package feed

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// Service implements FeedService.
type Service struct{}

// NewService creates a new feed service.
func NewService() *Service {
	return &Service{}
}

// ListFeed returns a chronological list of posts.
func (s *Service) ListFeed(ctx context.Context, req *connect.Request[stillv1.ListFeedRequest]) (*connect.Response[stillv1.ListFeedResponse], error) {
	return connect.NewResponse(&stillv1.ListFeedResponse{
		Posts: []*stillv1.Post{},
	}), nil
}
