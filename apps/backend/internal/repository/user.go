package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

// UserRepository handles user persistence.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// GetUser returns a user by ID.
func (r *UserRepository) GetUser(ctx context.Context, id string) (*stillv1.User, error) {
	var u stillv1.User
	var createdAt time.Time
	var avatarUrl *string
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, avatar_url, created_at FROM users WHERE id = $1
	`, id).Scan(&u.Id, &u.Username, &avatarUrl, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	if avatarUrl != nil {
		u.AvatarUrl = *avatarUrl
	}
	u.CreatedAt = timestamppb.New(createdAt)
	return &u, nil
}

// CountPosts returns the number of posts created by a user.
func (r *UserRepository) CountPosts(ctx context.Context, userID string) (int32, error) {
	var count int32
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count posts failed: %w", err)
	}
	return count, nil
}
