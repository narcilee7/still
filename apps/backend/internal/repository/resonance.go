package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ResonanceRepository handles resonance persistence.
type ResonanceRepository struct {
	pool *pgxpool.Pool
}

// NewResonanceRepository creates a new ResonanceRepository.
func NewResonanceRepository(pool *pgxpool.Pool) *ResonanceRepository {
	return &ResonanceRepository{pool: pool}
}

// ToggleResonance adds or removes a resonance for (postID, userID).
// It returns the delta (+1 or -1) applied to the post's resonance_count.
func (r *ResonanceRepository) ToggleResonance(ctx context.Context, postID, userID string) (int32, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM resonances WHERE post_id = $1 AND user_id = $2)
	`, postID, userID).Scan(&exists); err != nil {
		return 0, fmt.Errorf("check resonance failed: %w", err)
	}

	var delta int32
	if exists {
		_, err = tx.Exec(ctx, `
			DELETE FROM resonances WHERE post_id = $1 AND user_id = $2
		`, postID, userID)
		if err != nil {
			return 0, fmt.Errorf("delete resonance failed: %w", err)
		}
		delta = -1
	} else {
		_, err = tx.Exec(ctx, `
			INSERT INTO resonances (post_id, user_id) VALUES ($1, $2)
		`, postID, userID)
		if err != nil {
			return 0, fmt.Errorf("insert resonance failed: %w", err)
		}
		delta = 1
	}

	_, err = tx.Exec(ctx, `
		UPDATE posts SET resonance_count = GREATEST(0, resonance_count + $2) WHERE id = $1
	`, postID, delta)
	if err != nil {
		return 0, fmt.Errorf("update resonance_count failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit resonance failed: %w", err)
	}
	return delta, nil
}

// HasResonated returns true if the user has resonated with the post.
func (r *ResonanceRepository) HasResonated(ctx context.Context, postID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM resonances WHERE post_id = $1 AND user_id = $2)
	`, postID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("has resonated failed: %w", err)
	}
	return exists, nil
}

// CountByUser returns the total number of resonances received by a user across all their posts.
func (r *ResonanceRepository) CountByUser(ctx context.Context, userID string) (int32, error) {
	var count int32
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(p.resonance_count), 0)
		FROM posts p
		WHERE p.user_id = $1
	`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count user resonances failed: %w", err)
	}
	return count, nil
}
