package post

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/auth"
)

type fakePostStore struct {
	created bool
	post    *stillv1.Post
	err     error
}

func (f *fakePostStore) CreatePost(ctx context.Context, userID, imageURL, mood, title, description, status string) (*stillv1.Post, error) {
	f.created = true
	if f.post != nil {
		return f.post, nil
	}
	return &stillv1.Post{
		Id:       "post-1",
		UserId:   userID,
		ImageUrl: imageURL,
		Mood:     mood,
		Title:    title,
		Status:   stillv1.PostStatus_POST_STATUS_APPROVED,
	}, nil
}

func (f *fakePostStore) GetPost(ctx context.Context, id string) (*stillv1.Post, error) {
	if f.post != nil {
		return f.post, nil
	}
	return &stillv1.Post{Id: id}, nil
}

func (f *fakePostStore) UpdatePost(ctx context.Context, postID, userID, mood, title, description string) (*stillv1.Post, error) {
	return &stillv1.Post{
		Id:          postID,
		UserId:      userID,
		Mood:        mood,
		Title:       title,
		Description: description,
	}, nil
}

func (f *fakePostStore) DeletePost(ctx context.Context, postID, userID string) error {
	return nil
}

type fakeAnalysisStore struct {
	saved bool
	err   error
}

func (f *fakeAnalysisStore) SaveAnalysis(ctx context.Context, postID, provider, promptVersion string, rawResponse any) error {
	f.saved = true
	return f.err
}

type fakeUserReader struct {
	user *stillv1.User
	err  error
}

func (f *fakeUserReader) GetOrCreateUserByClerkID(ctx context.Context, clerkUserID, username string) (*stillv1.User, error) {
	if f.user != nil {
		return f.user, nil
	}
	return &stillv1.User{
		Id:        "user-1",
		Username:  "test_user",
		CreatedAt: timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}, f.err
}

type fakeAnalyzer struct {
	called bool
	result *ai.Result
	err    error
}

func (f *fakeAnalyzer) Analyze(ctx context.Context, imageURL string) (*ai.Result, error) {
	f.called = true
	return f.result, f.err
}

func TestCreatePost_SkipsAI_WhenAllFieldsProvided(t *testing.T) {
	postStore := &fakePostStore{}
	analysisStore := &fakeAnalysisStore{}
	analyzer := &fakeAnalyzer{}
	userReader := &fakeUserReader{}

	svc := &Service{
		postRepo:     postStore,
		analysisRepo: analysisStore,
		userRepo:     userReader,
		analyzer:     analyzer,
	}

	req := connect.NewRequest(&stillv1.CreatePostRequest{
		ImageUrl:    "https://example.com/photo.jpg",
		Mood:        "still",
		Title:       "A quiet title",
		Description: "A short description",
	})
	ctx := auth.WithUserID(context.Background(), "clerk_user_123")

	_, err := svc.CreatePost(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if analyzer.called {
		t.Error("analyzer should not be called when all fields are provided")
	}
	if analysisStore.saved {
		t.Error("analysis should not be saved when AI is not called")
	}
}

func TestCreatePost_CallsAI_WhenFieldsMissing(t *testing.T) {
	postStore := &fakePostStore{}
	analysisStore := &fakeAnalysisStore{}
	analyzer := &fakeAnalyzer{
		result: &ai.Result{
			Mood:        "still",
			Title:       "AI title",
			Description: "AI description",
		},
	}
	userReader := &fakeUserReader{}

	svc := &Service{
		postRepo:     postStore,
		analysisRepo: analysisStore,
		userRepo:     userReader,
		analyzer:     analyzer,
	}

	req := connect.NewRequest(&stillv1.CreatePostRequest{
		ImageUrl: "https://example.com/photo.jpg",
	})
	ctx := auth.WithUserID(context.Background(), "clerk_user_123")

	_, err := svc.CreatePost(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !analyzer.called {
		t.Error("analyzer should be called when fields are missing")
	}
	if !analysisStore.saved {
		t.Error("analysis should be saved when AI is called")
	}
}

func TestCreatePost_ValidatesEmptyImageUrl(t *testing.T) {
	svc := &Service{}
	req := connect.NewRequest(&stillv1.CreatePostRequest{})

	_, err := svc.CreatePost(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty image_url")
	}
	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", connect.CodeOf(err))
	}
}

func TestCreatePost_RequiresAuthentication(t *testing.T) {
	svc := &Service{
		postRepo:     &fakePostStore{},
		analysisRepo: &fakeAnalysisStore{},
		userRepo:     &fakeUserReader{},
		analyzer:     &fakeAnalyzer{},
	}
	req := connect.NewRequest(&stillv1.CreatePostRequest{
		ImageUrl: "https://example.com/photo.jpg",
		Mood:     "still",
		Title:    "A quiet title",
	})

	_, err := svc.CreatePost(context.Background(), req)
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Errorf("expected Unauthenticated, got %v", connect.CodeOf(err))
	}
}
