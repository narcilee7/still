package user

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/auth"
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

// GetMe returns the currently authenticated user.
func (s *Service) GetMe(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[stillv1.GetMeResponse], error) {
	user, err := auth.CurrentUser(ctx, s.userRepo)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&stillv1.GetMeResponse{User: user}), nil
}

// GetProfile returns a user's profile.
func (s *Service) GetProfile(ctx context.Context, req *connect.Request[stillv1.GetProfileRequest]) (*connect.Response[stillv1.GetProfileResponse], error) {
	if strings.TrimSpace(req.Msg.UserId) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user_id is required"))
	}
	user, err := s.userRepo.GetUser(ctx, req.Msg.UserId)
	if err != nil {
		code := connect.CodeInternal
		if strings.Contains(err.Error(), "not found") {
			code = connect.CodeNotFound
		}
		return nil, connect.NewError(code, err)
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
		User:            user,
		PostsCount:      postsCount,
		ResonancesCount: resonancesCount,
	}), nil
}
