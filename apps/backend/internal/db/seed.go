package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DefaultUserID is the seeded demo user used until real auth is implemented.
const DefaultUserID = "00000000-0000-0000-0000-000000000001"

// Seed inserts default user and sample posts if the database is empty.
func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	if err := seedUser(ctx, pool); err != nil {
		return err
	}
	return seedPosts(ctx, pool)
}

func seedUser(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, username, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (username) DO NOTHING
	`, DefaultUserID, "still_moments")
	return err
}

func seedPosts(ctx context.Context, pool *pgxpool.Pool) error {
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM posts`).Scan(&count); err != nil {
		return fmt.Errorf("count posts failed: %w", err)
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC()
	posts := []struct {
		imageURL        string
		mood            string
		title           string
		description     string
		resonanceCount  int32
		createdAtOffset time.Duration
	}{
		{"https://picsum.photos/seed/still001/800/1000", "still", "The Hour Before Rain", "Everything slows down when the sky turns grey.", 28, -time.Hour * 2},
		{"https://picsum.photos/seed/still002/800/1000", "drift", "Passing Through", "The city keeps moving, even when your thoughts do not.", 14, -time.Hour * 3},
		{"https://picsum.photos/seed/still003/800/1000", "warm", "Afternoon Light", "A quiet reminder that some things are still soft.", 42, -time.Hour * 24},
		{"https://picsum.photos/seed/still004/800/1000", "waiting", "On The Way", "Some days are made of unfinished thoughts.", 7, -time.Hour * 48},
		{"https://picsum.photos/seed/still005/800/1000", "distant", "Faraway Sounds", "Listening to a conversation that is not yours.", 19, -time.Hour * 72},
		{"https://picsum.photos/seed/still006/800/1000", "quiet", "Before Speaking", "The pause holds more than the words ever could.", 33, -time.Hour * 96},
		{"https://picsum.photos/seed/still007/800/1000", "soft", "Sunday Morning", "Light through curtains, coffee going cold.", 56, -time.Hour * 120},
		{"https://picsum.photos/seed/still008/800/1000", "returning", "Back Again", "The same street, a different version of you.", 12, -time.Hour * 144},
	}

	batch := &pgx.Batch{}
	for _, p := range posts {
		batch.Queue(`
			INSERT INTO posts (user_id, image_url, mood, title, description, status, resonance_count, created_at)
			VALUES ($1, $2, $3, $4, $5, 'APPROVED', $6, $7)
		`, DefaultUserID, p.imageURL, p.mood, p.title, p.description, p.resonanceCount, now.Add(p.createdAtOffset))
	}
	br := pool.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("seed posts failed: %w", err)
	}
	return nil
}
