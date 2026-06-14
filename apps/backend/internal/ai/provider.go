package ai

import "strings"

// Provider identifies the LLM backend.
type Provider string

const (
	ProviderOpenAI   Provider = "openai"
	ProviderDeepSeek Provider = "deepseek"
	ProviderMoonshot Provider = "moonshot"
	ProviderQwen     Provider = "qwen"
)

// ParseProvider normalizes a provider name.
func ParseProvider(s string) Provider {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "deepseek":
		return ProviderDeepSeek
	case "moonshot":
		return ProviderMoonshot
	case "qwen", "dashscope":
		return ProviderQwen
	default:
		return ProviderOpenAI
	}
}

// SupportsVision reports whether the provider natively supports image inputs.
//
// Notes:
//   - DeepSeek's V4 family supports vision, but only on the Pro model
//     (deepseek-v4-pro). The default is set to that model so image inputs work
//     out of the box; v4-flash is text-only.
func (p Provider) SupportsVision() bool {
	switch p {
	case ProviderOpenAI, ProviderQwen, ProviderMoonshot:
		return true
	case ProviderDeepSeek:
		return false
	default:
		return false
	}
}

func (p Provider) String() string { return string(p) }
