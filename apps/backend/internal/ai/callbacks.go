package ai

import (
	"context"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	cb "github.com/cloudwego/eino/utils/callbacks"
	"github.com/rs/zerolog/log"
)

func newUsageCallbackHandler() callbacks.Handler {
	return cb.NewHandlerHelper().ChatModel(&cb.ModelCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks.RunInfo, input *model.CallbackInput) context.Context {
			log.Ctx(ctx).Debug().
				Str("component", runInfo.Name).
				Str("model", input.Config.Model).
				Int("messages", len(input.Messages)).
				Msg("llm request started")
			return ctx
		},
		OnEnd: func(ctx context.Context, runInfo *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			event := log.Ctx(ctx).Debug().
				Str("component", runInfo.Name).
				Str("model", output.Config.Model)
			if output.TokenUsage != nil {
				event = event.
					Int("prompt_tokens", output.TokenUsage.PromptTokens).
					Int("completion_tokens", output.TokenUsage.CompletionTokens).
					Int("total_tokens", output.TokenUsage.TotalTokens)
			}
			event.Msg("llm request completed")
			return ctx
		},
		OnError: func(ctx context.Context, runInfo *callbacks.RunInfo, err error) context.Context {
			log.Ctx(ctx).Error().
				Err(err).
				Str("component", runInfo.Name).
				Msg("llm request failed")
			return ctx
		},
	}).Handler()
}
