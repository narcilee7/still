package post

import (
	"context"

	"connectrpc.com/connect"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/db"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// Service implements PostService.
type Service struct {
	postRepo    *repository.PostRepository
	analysisRepo *repository.AnalysisRepository
	analyzer    ai.Analyzer
}

// NewService creates a new post service.
func NewService(postRepo *repository.PostRepository, analysisRepo *repository.AnalysisRepository, analyzer ai.Analyzer) *Service {
	return &Service{
		postRepo:    postRepo,
		analysisRepo: analysisRepo,
		analyzer:    analyzer,
	}
}

// CreatePost creates a new post.
func (s *Service) CreatePost(ctx context.Context, req *connect.Request[stillv1.CreatePostRequest]) (*connect.Response[stillv1.CreatePostResponse], error) {
	result, err := s.analyzer.Analyze(ctx, req.Msg.ImageUrl)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Use the provided mood/title/description from the request (user is author),
	// but fall back to AI analysis if empty.
	mood := req.Msg.Mood
	if mood == "" {
		mood = result.Mood
	}
	title := req.Msg.Title
	if title == "" {
		title = result.Title
	}
	description := req.Msg.Description
	if description == "" {
		description = result.Description
	}

	status := repository.StatusToString(stillv1.PostStatus_POST_STATUS_APPROVED)
	post, err := s.postRepo.CreatePost(ctx, db.DefaultUserID, req.Msg.ImageUrl, mood, title, description, status)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	provider := "stub"
	if _, ok := s.analyzer.(*ai.OpenAIAnalyzer); ok {
		provider = "openai"
	}
	if err := s.analysisRepo.SaveAnalysis(ctx, post.Id, provider, "v1", result); err != nil {
		// Non-fatal for the response; log would go here.
		_ = err
	}

	return connect.NewResponse(&stillv1.CreatePostResponse{Post: post}), nil
}

// GetPost returns a post by ID.
func (s *Service) GetPost(ctx context.Context, req *connect.Request[stillv1.GetPostRequest]) (*connect.Response[stillv1.GetPostResponse], error) {
	post, err := s.postRepo.GetPost(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(&stillv1.GetPostResponse{Post: post}), nil
}
