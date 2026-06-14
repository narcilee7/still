package ai

import (
	"context"
	"fmt"

	einocallbacks "github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	einoopenai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/qwen"

	"github.com/still-mvp/still/apps/backend/pkg/config"
)

// NewAnalyzerFromConfig builds an Analyzer from the application config.
// It also registers global Eino callbacks for usage / trajectory logging.
func NewAnalyzerFromConfig(ctx context.Context, cfg *config.Config) (Analyzer, error) {
	RegisterUsageCallbacks()

	provider := ParseProvider(cfg.LLMProvider)

	var cm model.BaseChatModel
	var err error

	switch provider {
	case ProviderOpenAI:
		cm, err = einoopenai.NewChatModel(ctx, &einoopenai.ChatModelConfig{
			APIKey: cfg.LLMAPIKey,
			Model:  cfg.LLMModel,
			BaseURL: firstNonEmpty(cfg.LLMBaseURL, "https://api.openai.com/v1"),
		})
	case ProviderDeepSeek:
		cm, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			APIKey: cfg.LLMAPIKey,
			Model:  cfg.LLMModel,
			BaseURL: firstNonEmpty(cfg.LLMBaseURL, "https://api.deepseek.com/"),
		})
	case ProviderMoonshot:
		// Moonshot exposes an OpenAI-compatible API; no dedicated Eino component yet.
		cm, err = einoopenai.NewChatModel(ctx, &einoopenai.ChatModelConfig{
			APIKey: cfg.LLMAPIKey,
			Model:  cfg.LLMModel,
			BaseURL: firstNonEmpty(cfg.LLMBaseURL, "https://api.moonshot.cn/v1"),
		})
	case ProviderQwen:
		cm, err = qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			APIKey: cfg.LLMAPIKey,
			Model:  cfg.LLMModel,
			BaseURL: firstNonEmpty(cfg.LLMBaseURL, "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		})
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.LLMProvider)
	}
	if err != nil {
		return nil, fmt.Errorf("create %s chat model: %w", provider, err)
	}

	return NewEinoAnalyzer(cm, provider), nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

var usageCallbackRegistered bool

// RegisterUsageCallbacks registers a global Eino callback handler that logs
// token usage and model calls for observability. Safe to call multiple times.
func RegisterUsageCallbacks() {
	if usageCallbackRegistered {
		return
	}
	usageCallbackRegistered = true

	einocallbacks.AppendGlobalHandlers(newUsageCallbackHandler())
}
