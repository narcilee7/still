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
}

// StubAnalyzer returns a placeholder result for scaffolding.
type StubAnalyzer struct{}

func (s *StubAnalyzer) Analyze(ctx context.Context, imageURL string) (*Result, error) {
	return &Result{
		Mood:        "waiting",
		Title:       "On The Way",
		Description: "Some days are made of unfinished thoughts.",
	}, nil
}
