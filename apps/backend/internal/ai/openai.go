package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const moodList = "still, waiting, drift, warm, distant, quiet, hollow, soft, fading, wondering, returning"

const analyzePrompt = `You are a quiet, sensitive photo editor for a mood-sharing community called Still.
Look at the photo and return a single JSON object with exactly these keys:
- "mood": one word chosen ONLY from this list: ` + moodList + `.
- "title": a short, poetic title (max 8 words).
- "description": one gentle sentence describing the feeling (max 25 words).

Do not explain, do not add markdown, do not include confidence scores. Return valid JSON only.`

// OpenAIAnalyzer uses OpenAI vision models to analyze images.
type OpenAIAnalyzer struct {
	client *openai.Client
	model  string
}

// NewOpenAIAnalyzer creates a new OpenAI-backed analyzer.
func NewOpenAIAnalyzer(apiKey string) *OpenAIAnalyzer {
	return &OpenAIAnalyzer{
		client: openai.NewClient(apiKey),
		model:  openai.GPT4oMini,
	}
}

// Analyze calls OpenAI vision and returns an emotional interpretation.
func (a *OpenAIAnalyzer) Analyze(ctx context.Context, imageURL string) (*Result, error) {
	resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: analyzePrompt,
			},
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{Type: openai.ChatMessagePartTypeText, Text: "What mood does this photo hold?"},
					{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{URL: imageURL}},
				},
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
		MaxTokens:      256,
	})
	if err != nil {
		return nil, fmt.Errorf("openai completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	var parsed struct {
		Mood        string `json:"mood"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("parse openai response failed (%q): %w", content, err)
	}

	return &Result{
		Mood:        normalizeMood(parsed.Mood),
		Title:       strings.TrimSpace(parsed.Title),
		Description: strings.TrimSpace(parsed.Description),
	}, nil
}

func normalizeMood(m string) string {
	m = strings.ToLower(strings.TrimSpace(m))
	for _, valid := range strings.Split(moodList, ", ") {
		if m == valid {
			return m
		}
	}
	return "still"
}
