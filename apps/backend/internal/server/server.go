package server

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/still-mvp/still/apps/backend/gen/still/v1/stillv1connect"
	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/analyze"
	"github.com/still-mvp/still/apps/backend/internal/feed"
	"github.com/still-mvp/still/apps/backend/internal/post"
	"github.com/still-mvp/still/apps/backend/internal/resonate"
	"github.com/still-mvp/still/apps/backend/internal/user"
)

// Server wraps the HTTP server and dependencies.
type Server struct {
	http *http.Server
}

// New creates a new server.
func New(addr string) *Server {
	mux := http.NewServeMux()

	analyzer := &ai.StubAnalyzer{}

	mux.Handle(stillv1connect.NewFeedServiceHandler(feed.NewService()))
	mux.Handle(stillv1connect.NewPostServiceHandler(post.NewService()))
	mux.Handle(stillv1connect.NewAnalyzeServiceHandler(analyze.NewService(analyzer)))
	mux.Handle(stillv1connect.NewResonateServiceHandler(resonate.NewService()))
	mux.Handle(stillv1connect.NewUserServiceHandler(user.NewService()))

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: withCORS(mux),
		},
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	log.Info().Str("addr", s.http.Addr).Msg("starting server")
	go func() {
		<-ctx.Done()
		_ = s.http.Shutdown(context.Background())
	}()
	return s.http.ListenAndServe()
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, X-User-Agent")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
