package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// PostRepository handles post persistence.
type PostRepository struct {
	pool *pgxpool.Pool
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

// CreatePost inserts a new post and returns it.
func (r *PostRepository) CreatePost(ctx context.Context, userID, imageURL, mood, title, description, status string) (*stillv1.Post, error) {
	var id string
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, `
		INSERT INTO posts (user_id, image_url, mood, title, description, status, resonance_count)
		VALUES ($1, $2, $3, $4, $5, $6, 0)
		RETURNING id, created_at
	`, userID, imageURL, mood, title, description, status).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("create post failed: %w", err)
	}
	return &stillv1.Post{
		Id:             id,
		UserId:         userID,
		ImageUrl:       imageURL,
		Mood:           mood,
		Title:          title,
		Description:    description,
		Status:         parseStatus(status),
		CreatedAt:      timestamppb.New(createdAt),
		ResonanceCount: 0,
	}, nil
}

// GetPost returns a post by ID.
func (r *PostRepository) GetPost(ctx context.Context, id string) (*stillv1.Post, error) {
	post, err := r.scanPost(ctx, `
		SELECT id, user_id, image_url, mood, title, description, status, created_at, resonance_count
		FROM posts
		WHERE id = $1
	`, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}
	return post, nil
}

// ListPosts returns posts ordered by created_at DESC with keyset pagination.
// pageToken is an RFC3339 timestamp; only posts older than that token are returned.
func (r *PostRepository) ListPosts(ctx context.Context, pageSize int32, pageToken string) ([]*stillv1.Post, string, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	args := []any{pageSize}
	where := ""
	if pageToken != "" {
		t, err := time.Parse(time.RFC3339Nano, pageToken)
		if err != nil {
			return nil, "", fmt.Errorf("invalid page_token: %w", err)
		}
		where = "WHERE created_at < $2"
		args = append(args, t)
	}

	statusWhere := "status = 'APPROVED'"
	if where != "" {
		statusWhere = where + " AND " + statusWhere
	} else {
		statusWhere = "WHERE " + statusWhere
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, image_url, mood, title, description, status, created_at, resonance_count
		FROM posts
		`+statusWhere+`
		ORDER BY created_at DESC
		LIMIT $1
	`, args...)
	if err != nil {
		return nil, "", fmt.Errorf("list posts failed: %w", err)
	}
	defer rows.Close()

	posts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*stillv1.Post, error) {
		return scanPostRow(row)
	})
	if err != nil {
		return nil, "", fmt.Errorf("scan posts failed: %w", err)
	}

	nextToken := ""
	if len(posts) == int(pageSize) {
		nextToken = posts[len(posts)-1].CreatedAt.AsTime().Format(time.RFC3339Nano)
	}
	return posts, nextToken, nil
}

// UpdateResonanceCount sets the resonance_count for a post.
func (r *PostRepository) UpdateResonanceCount(ctx context.Context, postID string, delta int32) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE posts
		SET resonance_count = GREATEST(0, resonance_count + $2)
		WHERE id = $1
	`, postID, delta)
	return err
}

// UpdatePost updates a post's mood, title and description if it belongs to userID.
func (r *PostRepository) UpdatePost(ctx context.Context, postID, userID, mood, title, description string) (*stillv1.Post, error) {
	var id string
	var createdAt time.Time
	var imageURL string
	var status string
	var resonanceCount int32
	err := r.pool.QueryRow(ctx, `
		UPDATE posts
		SET mood = $3, title = $4, description = $5
		WHERE id = $1 AND user_id = $2
		RETURNING id, image_url, status, created_at, resonance_count
	`, postID, userID, mood, title, description).Scan(&id, &imageURL, &status, &createdAt, &resonanceCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found or not owned by user")
		}
		return nil, fmt.Errorf("update post failed: %w", err)
	}
	return &stillv1.Post{
		Id:             id,
		UserId:         userID,
		ImageUrl:       imageURL,
		Mood:           mood,
		Title:          title,
		Description:    description,
		Status:         parseStatus(status),
		CreatedAt:      timestamppb.New(createdAt),
		ResonanceCount: resonanceCount,
	}, nil
}

// DeletePost removes a post if it belongs to userID.
func (r *PostRepository) DeletePost(ctx context.Context, postID, userID string) error {
	res, err := r.pool.Exec(ctx, `DELETE FROM posts WHERE id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		return fmt.Errorf("delete post failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("post not found or not owned by user")
	}
	return nil
}

func (r *PostRepository) scanPost(ctx context.Context, sql string, args ...any) (*stillv1.Post, error) {
	row := r.pool.QueryRow(ctx, sql, args...)
	return scanPostRow(row)
}

func scanPostRow(row pgx.Row) (*stillv1.Post, error) {
	var p stillv1.Post
	var createdAt time.Time
	var status string
	err := row.Scan(&p.Id, &p.UserId, &p.ImageUrl, &p.Mood, &p.Title, &p.Description, &status, &createdAt, &p.ResonanceCount)
	if err != nil {
		return nil, err
	}
	p.Status = parseStatus(status)
	p.CreatedAt = timestamppb.New(createdAt)
	return &p, nil
}

func parseStatus(s string) stillv1.PostStatus {
	switch strings.ToUpper(s) {
	case "PENDING":
		return stillv1.PostStatus_POST_STATUS_PENDING
	case "APPROVED":
		return stillv1.PostStatus_POST_STATUS_APPROVED
	case "REJECTED":
		return stillv1.PostStatus_POST_STATUS_REJECTED
	default:
		return stillv1.PostStatus_POST_STATUS_UNSPECIFIED
	}
}

// StatusToString converts a proto status to its DB string.
func StatusToString(s stillv1.PostStatus) string {
	switch s {
	case stillv1.PostStatus_POST_STATUS_PENDING:
		return "PENDING"
	case stillv1.PostStatus_POST_STATUS_APPROVED:
		return "APPROVED"
	case stillv1.PostStatus_POST_STATUS_REJECTED:
		return "REJECTED"
	default:
		return "PENDING"
	}
}
