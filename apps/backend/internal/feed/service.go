package feed

import (
	"context"
	"strings"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// Service implements FeedService.
type Service struct {
	postRepo *repository.PostRepository
}

// NewService creates a new feed service.
func NewService(postRepo *repository.PostRepository) *Service {
	return &Service{postRepo: postRepo}
}

// ListFeed returns a chronological list of posts.
func (s *Service) ListFeed(ctx context.Context, req *connect.Request[stillv1.ListFeedRequest]) (*connect.Response[stillv1.ListFeedResponse], error) {
	posts, nextPageToken, err := s.postRepo.ListPosts(ctx, req.Msg.PageSize, req.Msg.PageToken)
	if err != nil {
		code := connect.CodeInternal
		if strings.Contains(err.Error(), "invalid page_token") {
			code = connect.CodeInvalidArgument
		}
		return nil, connect.NewError(code, err)
	}
	return connect.NewResponse(&stillv1.ListFeedResponse{
		Posts:         posts,
		NextPageToken: nextPageToken,
	}), nil
}
