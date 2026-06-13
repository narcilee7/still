package user

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// Service implements UserService.
type Service struct {
	userRepo      *repository.UserRepository
	resonanceRepo *repository.ResonanceRepository
}

// NewService creates a new user service.
func NewService(userRepo *repository.UserRepository, resonanceRepo *repository.ResonanceRepository) *Service {
	return &Service{
		userRepo:      userRepo,
		resonanceRepo: resonanceRepo,
	}
}

// GetProfile returns a user's profile.
func (s *Service) GetProfile(ctx context.Context, req *connect.Request[stillv1.GetProfileRequest]) (*connect.Response[stillv1.GetProfileResponse], error) {
	user, err := s.userRepo.GetUser(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	postsCount, err := s.userRepo.CountPosts(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resonancesCount, err := s.resonanceRepo.CountByUser(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&stillv1.GetProfileResponse{
		User:             user,
		PostsCount:       postsCount,
		ResonancesCount:  resonancesCount,
	}), nil
}
