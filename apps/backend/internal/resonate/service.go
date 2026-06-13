package resonate

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/db"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// Service implements ResonateService.
type Service struct {
	postRepo      *repository.PostRepository
	resonanceRepo *repository.ResonanceRepository
}

// NewService creates a new resonate service.
func NewService(postRepo *repository.PostRepository, resonanceRepo *repository.ResonanceRepository) *Service {
	return &Service{
		postRepo:      postRepo,
		resonanceRepo: resonanceRepo,
	}
}

// Resonate records a resonance on a post.
func (s *Service) Resonate(ctx context.Context, req *connect.Request[stillv1.ResonateRequest]) (*connect.Response[stillv1.ResonateResponse], error) {
	if _, err := s.resonanceRepo.ToggleResonance(ctx, req.Msg.PostId, db.DefaultUserID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	post, err := s.postRepo.GetPost(ctx, req.Msg.PostId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&stillv1.ResonateResponse{Post: post}), nil
}
