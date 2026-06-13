package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/still-mvp/still/apps/backend/gen/still/v1/stillv1connect"
	"github.com/still-mvp/still/apps/backend/internal/ai"
	"github.com/still-mvp/still/apps/backend/internal/analyze"
	"github.com/still-mvp/still/apps/backend/internal/feed"
	"github.com/still-mvp/still/apps/backend/internal/post"
	"github.com/still-mvp/still/apps/backend/internal/repository"
	"github.com/still-mvp/still/apps/backend/internal/resonate"
	"github.com/still-mvp/still/apps/backend/internal/storage"
	"github.com/still-mvp/still/apps/backend/internal/user"
	"github.com/still-mvp/still/apps/backend/pkg/config"
)

// Server wraps the HTTP server and dependencies.
type Server struct {
	http *http.Server
}

// New creates a new server.
func New(addr string, pool *pgxpool.Pool, analyzer ai.Analyzer, store storage.Store, cfg *config.Config) *Server {
	mux := http.NewServeMux()

	postRepo := repository.NewPostRepository(pool)
	resonanceRepo := repository.NewResonanceRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	analysisRepo := repository.NewAnalysisRepository(pool)

	mux.Handle(stillv1connect.NewFeedServiceHandler(feed.NewService(postRepo)))
	mux.Handle(stillv1connect.NewPostServiceHandler(post.NewService(postRepo, analysisRepo, analyzer)))
	mux.Handle(stillv1connect.NewAnalyzeServiceHandler(analyze.NewService(analyzer)))
	mux.Handle(stillv1connect.NewResonateServiceHandler(resonate.NewService(postRepo, resonanceRepo)))
	mux.Handle(stillv1connect.NewUserServiceHandler(user.NewService(userRepo, resonanceRepo)))
	mux.Handle(stillv1connect.NewStorageServiceHandler(storage.NewService(store)))

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: withCORS(mux, cfg),
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

func withCORS(next http.Handler, cfg *config.Config) http.Handler {
	allowed := cfg.CORSAllowedOrigins
	if allowed == "" {
		allowed = "*"
	}
	origins := strings.Split(allowed, ",")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && contains(origins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, X-User-Agent, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}
