package webserver

import (
	"context"
	"net/http"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
)

// Router represents HTTP route multiplexer
type Router struct {
	*chi.Mux
}

// Server - serves http requests
type Server struct {
	server *http.Server
}

// NewRouter - returns new router
func NewRouter(logger middleware.LogFormatter) *Router {
	router := &Router{chi.NewRouter()}
	router.Use(chiprometheus.NewMiddleware("tugc-aggregator", 5, 50, 100, 500, 1000, 2000))
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Recoverer)
	router.Mount("/debug", middleware.Profiler())
	return router
}

func NewCommonRouter() *Router {
	router := &Router{chi.NewRouter()}
	router.Use(middleware.Recoverer)
	router.Mount("/debug", middleware.Profiler())
	return router
}

// NewServer инициализация севрера с настройками и зависимостями
func NewServer(addr string, h http.Handler) *Server {
	return &Server{server: &http.Server{Addr: addr, Handler: h}}
}

// Start http server
func (s *Server) Start() error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.server.ListenAndServe()
	}()
	if err := <-errChan; err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "server error")
		}
	}
	return nil
}

// Stop http server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "could not shutdown tugc-aggregator server")
	}
	return nil
}
