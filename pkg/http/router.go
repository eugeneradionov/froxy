package http

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) MountRoutes(proxyTr http.Handler) {
	s.router.Use(middleware.Recoverer, middleware.RequestID)

	s.router.HandleFunc("/health", health)

	s.router.Mount("/proxy/v1", proxyTr)
}
