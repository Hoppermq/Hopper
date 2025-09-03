package ui

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/hopper/internal/ui/routes"
	"github.com/hoppermq/hopper/pkg/domain"
)

type HTTPServer struct {
	logger *slog.Logger
	engine *gin.Engine
	server *http.Server
}

type Option func(*HTTPServer)

func WithEngine(engine *gin.Engine) Option {
	return func(s *HTTPServer) {
		s.engine = engine
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *HTTPServer) {
		s.logger = logger
	}
}

func NewHTTPServer(opts ...Option) *HTTPServer {
	httpServer := &HTTPServer{}

	for _, opt := range opts {
		opt(httpServer)
	}

	httpServer.server = &http.Server{
		Addr:    ":8090",
		Handler: httpServer.engine,
	}

	return httpServer
}

// Run start the ui http server.
func (s *HTTPServer) Run(ctx context.Context) error {
	routes.RegisterBaseRoutes(s.engine)
	if err := s.engine.Run(":8090"); err != nil && err != http.ErrServerClosed {
		s.logger.Info("http server closed", "error", err)
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Info("http server closed", "error", err)
		}
	}()

	<-ctx.Done()
	return nil
}

// Stop shutdown gracefully the ui http server.
func (s *HTTPServer) Stop(ctx context.Context) error {
	return nil
}

// Name return the service name.
func (s *HTTPServer) Name() string {
	return "hopper-ui"
}

func (s *HTTPServer) RegisterEventBus(eb domain.IEventBus) {
	// service not eb aware
}
