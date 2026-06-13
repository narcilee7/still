package post

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/auth"
	"github.com/still-mvp/still/apps/backend/internal/repository"
)

// PostStore creates and reads posts.
type PostStore interface {
	CreatePost(ctx context.Context, userID, imageURL, mood, title, description, status string) (*stillv1.Post, error)
	GetPost(ctx context.Context, id string) (*stillv1.Post, error)
}

// AnalysisStore persists raw AI analysis results.
type AnalysisStore interface {
	SaveAnalysis(ctx context.Context, postID, provider, promptVersion string, rawResponse any) error
}

// Service implements PostService.
type Service struct {
	postRepo     PostStore
	analysisRepo AnalysisStore
	userRepo     UserReader
	analyzer     ai.Analyzer
}

// UserReader resolves a Clerk identity to an internal user.
type UserReader interface {
	GetOrCreateUserByClerkID(ctx context.Context, clerkUserID, username string) (*stillv1.User, error)
}

// NewService creates a new post service.
func NewService(postRepo *repository.PostRepository, analysisRepo *repository.AnalysisRepository, userRepo UserReader, analyzer ai.Analyzer) *Service {
	return &Service{
		postRepo:     postRepo,
		analysisRepo: analysisRepo,
		userRepo:     userRepo,
		analyzer:     analyzer,
	}
}

// CreatePost creates a new post.
func (s *Service) CreatePost(ctx context.Context, req *connect.Request[stillv1.CreatePostRequest]) (*connect.Response[stillv1.CreatePostResponse], error) {
	if strings.TrimSpace(req.Msg.ImageUrl) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("image_url is required"))
	}

	needsAnalysis := req.Msg.Mood == "" || req.Msg.Title == "" || req.Msg.Description == ""

	var result *ai.Result
	var err error
	if needsAnalysis {
		result, err = s.analyzer.Analyze(ctx, req.Msg.ImageUrl)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Use the provided mood/title/description from the request (user is author),
	// but fall back to AI analysis if empty.
	mood := req.Msg.Mood
	if mood == "" && result != nil {
		mood = result.Mood
	}
	title := req.Msg.Title
	if title == "" && result != nil {
		title = result.Title
	}
	description := req.Msg.Description
	if description == "" && result != nil {
		description = result.Description
	}

	user, err := auth.CurrentUser(ctx, s.userRepo)
	if err != nil {
		return nil, err
	}

	status := repository.StatusToString(stillv1.PostStatus_POST_STATUS_APPROVED)
	post, err := s.postRepo.CreatePost(ctx, user.Id, req.Msg.ImageUrl, mood, title, description, status)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if result != nil {
		provider := "stub"
		if _, ok := s.analyzer.(*ai.OpenAIAnalyzer); ok {
			provider = "openai"
		}
		if err := s.analysisRepo.SaveAnalysis(ctx, post.Id, provider, "v1", result); err != nil {
			log.Warn().Err(err).Str("post_id", post.Id).Msg("save analysis failed")
		}
	}

	return connect.NewResponse(&stillv1.CreatePostResponse{Post: post}), nil
}

// GetPost returns a post by ID.
func (s *Service) GetPost(ctx context.Context, req *connect.Request[stillv1.GetPostRequest]) (*connect.Response[stillv1.GetPostResponse], error) {
	if strings.TrimSpace(req.Msg.Id) == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}
	post, err := s.postRepo.GetPost(ctx, req.Msg.Id)
	if err != nil {
		code := connect.CodeInternal
		if strings.Contains(err.Error(), "not found") {
			code = connect.CodeNotFound
		}
		return nil, connect.NewError(code, err)
	}
	return connect.NewResponse(&stillv1.GetPostResponse{Post: post}), nil
}
