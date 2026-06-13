package resonate

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// Service implements ResonateService.
type Service struct{}

// NewService creates a new resonate service.
func NewService() *Service {
	return &Service{}
}

// Resonate records a resonance on a post.
func (s *Service) Resonate(ctx context.Context, req *connect.Request[stillv1.ResonateRequest]) (*connect.Response[stillv1.ResonateResponse], error) {
	return connect.NewResponse(&stillv1.ResonateResponse{}), nil
}
