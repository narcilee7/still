package ai

import "context"

// Result is the output of an emotion analysis.
type Result struct {
	Mood        string
	Title       string
	Description string
}

// Analyzer analyzes an image and returns an emotional interpretation.
type Analyzer interface {
	Analyze(ctx context.Context, imageURL string) (*Result, error)
	// Provider returns the LLM provider name used for telemetry / storage.
	Provider() string
}
