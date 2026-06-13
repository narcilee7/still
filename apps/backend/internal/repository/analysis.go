package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AnalysisRepository records raw AI analysis results.
type AnalysisRepository struct {
	pool *pgxpool.Pool
}

// NewAnalysisRepository creates a new AnalysisRepository.
func NewAnalysisRepository(pool *pgxpool.Pool) *AnalysisRepository {
	return &AnalysisRepository{pool: pool}
}

// SaveAnalysis persists an AI analysis record.
func (r *AnalysisRepository) SaveAnalysis(ctx context.Context, postID, provider, promptVersion string, rawResponse any) error {
	rawJSON, err := json.Marshal(rawResponse)
	if err != nil {
		return fmt.Errorf("marshal raw response failed: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO post_ai_analysis (post_id, provider, prompt_version, raw_response)
		VALUES ($1, $2, $3, $4)
	`, postID, provider, promptVersion, rawJSON)
	if err != nil {
		return fmt.Errorf("save analysis failed: %w", err)
	}
	return nil
}
