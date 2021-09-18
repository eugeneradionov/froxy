package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/eugeneradionov/froxy/config"
	"github.com/eugeneradionov/froxy/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	log    logger.Logger
	router chi.Router
}

func NewServer(cfg config.HTTPServer, log logger.Logger) *Server {
	router := chi.NewRouter()
	server := &http.Server{
		Handler: router,
		Addr:    cfg.ListenURL,
	}

	srv := Server{
		router: router,
		server: server,
		log:    log,
	}

	return &srv
}

func (s *Server) Mount(path string, handler http.Handler) {
	s.router.Mount(path, handler)
}

func (s *Server) Serve(ctx context.Context) error {
	go s.handleCtxCancel(ctx)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleCtxCancel(ctx context.Context) {
	<-ctx.Done()

	s.log.Info("http server is shutting down")

	ctx10s, _ := context.WithTimeout(context.Background(), 10*time.Second) // nolint
	if err := s.server.Shutdown(ctx10s); err != nil {
		s.log.Error("failed to shutdown the server", zap.Error(err))
	}

	s.log.Info("http server stopped")
}
