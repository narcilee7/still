package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const moodList = "still, waiting, drift, warm, distant, quiet, hollow, soft, fading, wondering, returning"

const analyzePrompt = `You are a quiet, sensitive photo editor for a mood-sharing community called Still.
Look at the photo and return a single JSON object with exactly these keys:
- "mood": one word chosen ONLY from this list: ` + moodList + `.
- "title": a short, poetic title (max 8 words).
- "description": one gentle sentence describing the feeling (max 25 words).

Do not explain, do not add markdown, do not include confidence scores. Return valid JSON only.`

// EinoAnalyzer implements Analyzer using a CloudWeGo Eino ChatModel.
type EinoAnalyzer struct {
	model    model.BaseChatModel
	provider Provider
}

// NewEinoAnalyzer creates an analyzer backed by an Eino ChatModel.
func NewEinoAnalyzer(cm model.BaseChatModel, provider Provider) *EinoAnalyzer {
	return &EinoAnalyzer{model: cm, provider: provider}
}

// Provider returns the configured provider name.
func (a *EinoAnalyzer) Provider() string { return a.provider.String() }

// Analyze invokes the model and parses the JSON response.
func (a *EinoAnalyzer) Analyze(ctx context.Context, imageURL string) (*Result, error) {
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: analyzePrompt,
		},
		{
			Role: schema.User,
			UserInputMultiContent: []schema.MessageInputPart{
				{Type: schema.ChatMessagePartTypeText, Text: "What mood does this photo hold?"},
				{
					Type: schema.ChatMessagePartTypeImageURL,
					Image: &schema.MessageInputImage{
						MessagePartCommon: schema.MessagePartCommon{
							URL: toPtr(imageURL),
						},
					},
				},
			},
		},
	}

	resp, err := a.model.Generate(ctx, messages,
		model.WithTemperature(0.4),
		model.WithMaxTokens(256),
	)
	if err != nil {
		return nil, fmt.Errorf("model generate failed: %w", err)
	}

	content := strings.TrimSpace(resp.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var parsed struct {
		Mood        string `json:"mood"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("parse model response failed (%q): %w", content, err)
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

func toPtr[T any](v T) *T { return &v }
