package resonate

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/auth"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// Service implements ResonateService.
type Service struct {
	postRepo      *repository.PostRepository
	resonanceRepo *repository.ResonanceRepository
	userRepo      UserReader
}

// UserReader resolves a Clerk identity to an internal user.
type UserReader interface {
	GetOrCreateUserByClerkID(ctx context.Context, clerkUserID, username string) (*stillv1.User, error)
}

// NewService creates a new resonate service.
func NewService(postRepo *repository.PostRepository, resonanceRepo *repository.ResonanceRepository, userRepo UserReader) *Service {
	return &Service{
		postRepo:      postRepo,
		resonanceRepo: resonanceRepo,
		userRepo:      userRepo,
	}
}

// Resonate records a resonance on a post.
func (s *Service) Resonate(ctx context.Context, req *connect.Request[stillv1.ResonateRequest]) (*connect.Response[stillv1.ResonateResponse], error) {
	if strings.TrimSpace(req.Msg.PostId) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("post_id is required"))
	}
	user, err := auth.CurrentUser(ctx, s.userRepo)
	if err != nil {
		return nil, err
	}
	_, hasResonated, err := s.resonanceRepo.ToggleResonance(ctx, req.Msg.PostId, user.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	post, err := s.postRepo.GetPost(ctx, req.Msg.PostId)
	if err != nil {
		code := connect.CodeInternal
		if strings.Contains(err.Error(), "not found") {
			code = connect.CodeNotFound
		}
		return nil, connect.NewError(code, err)
	}
	return connect.NewResponse(&stillv1.ResonateResponse{Post: post, HasResonated: hasResonated}), nil
}
