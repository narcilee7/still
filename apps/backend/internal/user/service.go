package user

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// Service implements UserService.
type Service struct{}

// NewService creates a new user service.
func NewService() *Service {
	return &Service{}
}

// GetProfile returns a user's profile.
func (s *Service) GetProfile(ctx context.Context, req *connect.Request[stillv1.GetProfileRequest]) (*connect.Response[stillv1.GetProfileResponse], error) {
	return connect.NewResponse(&stillv1.GetProfileResponse{}), nil
}
