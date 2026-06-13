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

// GetUserByClerkID returns a user by their Clerk external ID.
func (r *UserRepository) GetUserByClerkID(ctx context.Context, clerkUserID string) (*stillv1.User, error) {
	var u stillv1.User
	var createdAt time.Time
	var avatarUrl *string
	err := r.pool.QueryRow(ctx, `
		SELECT id, username, avatar_url, created_at FROM users WHERE clerk_user_id = $1
	`, clerkUserID).Scan(&u.Id, &u.Username, &avatarUrl, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by clerk id failed: %w", err)
	}
	if avatarUrl != nil {
		u.AvatarUrl = *avatarUrl
	}
	u.CreatedAt = timestamppb.New(createdAt)
	return &u, nil
}

// GetOrCreateUserByClerkID returns an existing user or creates one from the
// Clerk identity. It is idempotent.
func (r *UserRepository) GetOrCreateUserByClerkID(ctx context.Context, clerkUserID, username string) (*stillv1.User, error) {
	if u, err := r.GetUserByClerkID(ctx, clerkUserID); err == nil {
		return u, nil
	}

	candidate := username
	if candidate == "" {
		candidate = "still_" + clerkUserID
	}

	var u stillv1.User
	var createdAt time.Time
	var avatarUrl *string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (username, clerk_user_id)
		VALUES (
			CASE
				WHEN EXISTS (SELECT 1 FROM users WHERE username = $1) THEN $1 || '_' || substr(md5(random()::text), 1, 6)
				ELSE $1
			END,
			$2
		)
		RETURNING id, username, avatar_url, created_at
	`, candidate, clerkUserID).Scan(&u.Id, &u.Username, &avatarUrl, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
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
