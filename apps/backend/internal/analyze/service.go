package analyze

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/ai"
)

// Service implements AnalyzeService.
type Service struct {
	analyzer ai.Analyzer
}

// NewService creates a new analyze service.
func NewService(analyzer ai.Analyzer) *Service {
	return &Service{analyzer: analyzer}
}

// AnalyzeImage analyzes an image and returns mood/title/description.
func (s *Service) AnalyzeImage(ctx context.Context, req *connect.Request[stillv1.AnalyzeImageRequest]) (*connect.Response[stillv1.AnalyzeImageResponse], error) {
	result, err := s.analyzer.Analyze(ctx, req.Msg.ImageUrl)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&stillv1.AnalyzeImageResponse{
		Mood:        result.Mood,
		Title:       result.Title,
		Description: result.Description,
	}), nil
}
