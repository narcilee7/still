package ai

import (
	"context"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	cb "github.com/cloudwego/eino/utils/callbacks"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// callbackState carries the OpenTelemetry span and timing between OnStart and
// OnEnd/OnError. It is stored in the callback context.
type callbackState struct {
	span      trace.Span
	startedAt time.Time
}

type callbackCtxKey struct{}

func newUsageCallbackHandler() callbacks.Handler {
	return cb.NewHandlerHelper().ChatModel(&cb.ModelCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks.RunInfo, input *model.CallbackInput) context.Context {
			startedAt := time.Now()
			ctx, span := otel.Tracer("still.ai").Start(ctx, "llm.generate",
				trace.WithAttributes(
					attribute.String("llm.component", runInfo.Name),
					attribute.String("llm.model", input.Config.Model),
					attribute.Int("llm.input_messages", len(input.Messages)),
				),
			)

			ctx = context.WithValue(ctx, callbackCtxKey{}, &callbackState{
				span:      span,
				startedAt: startedAt,
			})

			textCount, imageCount := countMessageParts(input.Messages)
			log.Ctx(ctx).Debug().
				Str("component", runInfo.Name).
				Str("model", input.Config.Model).
				Int("messages", len(input.Messages)).
				Int("text_parts", textCount).
				Int("image_parts", imageCount).
				Msg("llm request started")

			return ctx
		},
		OnEnd: func(ctx context.Context, runInfo *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			state, ok := ctx.Value(callbackCtxKey{}).(*callbackState)
			if !ok {
				return ctx
			}

			duration := time.Since(state.startedAt)
			state.span.SetAttributes(
				attribute.Float64("llm.duration_ms", float64(duration.Milliseconds())),
			)

			event := log.Ctx(ctx).Info().
				Str("component", runInfo.Name).
				Str("model", output.Config.Model).
				Dur("duration", duration)

			bc := &sentry.Breadcrumb{
				Category: "llm",
				Message:  "llm request completed",
				Data: map[string]interface{}{
					"component":   runInfo.Name,
					"model":       output.Config.Model,
					"duration_ms": duration.Milliseconds(),
				},
				Level: sentry.LevelInfo,
			}

			if output.TokenUsage != nil {
				event = event.
					Int("prompt_tokens", output.TokenUsage.PromptTokens).
					Int("completion_tokens", output.TokenUsage.CompletionTokens).
					Int("total_tokens", output.TokenUsage.TotalTokens)

				state.span.SetAttributes(
					attribute.Int("llm.prompt_tokens", output.TokenUsage.PromptTokens),
					attribute.Int("llm.completion_tokens", output.TokenUsage.CompletionTokens),
					attribute.Int("llm.total_tokens", output.TokenUsage.TotalTokens),
				)

				bc.Data["prompt_tokens"] = output.TokenUsage.PromptTokens
				bc.Data["completion_tokens"] = output.TokenUsage.CompletionTokens
				bc.Data["total_tokens"] = output.TokenUsage.TotalTokens
			}

			// Trajectory: record a preview of the model response at DEBUG only.
			if output.Message != nil && output.Message.Content != "" {
				event = event.Str("response_preview", truncate(output.Message.Content, 200))
			}

			event.Msg("llm request completed")
			state.span.SetStatus(codes.Ok, "")
			state.span.End()

			sentry.AddBreadcrumb(bc)
			return ctx
		},
		OnError: func(ctx context.Context, runInfo *callbacks.RunInfo, err error) context.Context {
			state, ok := ctx.Value(callbackCtxKey{}).(*callbackState)
			if ok {
				duration := time.Since(state.startedAt)
				state.span.RecordError(err)
				state.span.SetStatus(codes.Error, err.Error())
				state.span.SetAttributes(
					attribute.Float64("llm.duration_ms", float64(duration.Milliseconds())),
				)
				state.span.End()
			}

			log.Ctx(ctx).Error().
				Err(err).
				Str("component", runInfo.Name).
				Msg("llm request failed")

			sentry.CaptureException(err)
			return ctx
		},
	}).Handler()
}

func countMessageParts(msgs []*schema.Message) (textCount, imageCount int) {
	for _, m := range msgs {
		parts := m.UserInputMultiContent
		if len(parts) == 0 {
			if m.Content != "" {
				textCount++
			}
			continue
		}
		for _, p := range parts {
			switch p.Type {
			case schema.ChatMessagePartTypeText:
				textCount++
			case schema.ChatMessagePartTypeImageURL:
				imageCount++
			}
		}
	}
	return
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
