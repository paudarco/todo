package server

import (
	"context"
	"net/http"
	"time"

	"github.com/paudarco/todo/internal/config"
)

// Абстракция для работы с http-сервером
type Server struct {
	server *http.Server
}

// / Run запускает сервер с указанными в config параметрами
func (s *Server) Run(cfg config.Server, handler http.Handler) error {
	s.server = &http.Server{
		Addr:              cfg.Host + ":" + cfg.Port,
		MaxHeaderBytes:    1 << 20, // 1 MB
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
	}

	return s.server.ListenAndServe()
}

// Graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
