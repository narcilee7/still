package server

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

	localStore, ok := store.(*storage.LocalStore)
	if ok {
		mux.Handle("/_upload/", uploadHandler(localStore.UploadDir()))
		mux.Handle("/_files/", fileHandler(localStore.UploadDir()))
	}

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	_ = cfg

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

func uploadHandler(uploadDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		key := filepath.Base(r.URL.Path)
		if key == "" || key == "." || key == "/" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := os.MkdirAll(uploadDir, 0o755); err != nil {
			log.Error().Err(err).Msg("create upload dir failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		path := filepath.Join(uploadDir, key)
		file, err := os.Create(path)
		if err != nil {
			log.Error().Err(err).Msg("create upload file failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		const maxUploadSize = 20 << 20 // 20 MB
		if _, err := io.Copy(file, io.LimitReader(r.Body, maxUploadSize)); err != nil {
			log.Error().Err(err).Msg("save upload file failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func fileHandler(uploadDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		key := filepath.Base(r.URL.Path)
		path := filepath.Join(uploadDir, key)
		http.ServeFile(w, r, path)
	})
}
