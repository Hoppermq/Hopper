// Package http  represent the http service.
package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/http/routes"
	"github.com/hoppermq/hopper/pkg/domain"
)

// HTTP is the HTTP Server of hoppermq
type HTTP struct {
	config *config.Configuration
	logger *slog.Logger
	eb     domain.IEventBus
	cancel context.CancelFunc

	// should be a domain interface here
	engine *gin.Engine
	server *http.Server
}

// Option is the type that represent the function to configure the server.
type Option func(*HTTP)

// WithLogger inject the logger to the http server.
func WithLogger(logger *slog.Logger) Option {
	return func(h *HTTP) {
		h.logger = logger
	}
}

// WithContext inject the cancelfunc ctx.
func WithContext(ctx context.CancelFunc) Option {
	return func(h *HTTP) {
		h.cancel = ctx
	}
}

// WithConfiguration inject the configuration.
func WithConfiguration(cfg *config.Configuration) Option {
	return func(h *HTTP) {
		h.config = cfg
	}
}

// WithEngine inject the http engine.
func WithEngine(engine *gin.Engine) Option {
	return func(h *HTTP) {
		h.engine = engine
	}
}

// NewHTTPServer return a new HTTP.
func NewHTTPServer(opts ...Option) *HTTP {
	httpServer := &HTTP{}

	for _, opt := range opts {
		opt(httpServer)
	}

	httpServer.server = &http.Server{
		Addr:         ":8080",
		Handler:      httpServer.engine,
		ReadTimeout:  10,
		WriteTimeout: 10,
	}

	return httpServer
}

// Run start the services.
func (h *HTTP) Run(ctx context.Context) error {
	h.logger.Info("starting the http server component", "name", h.Name())
	h.eb.Subscribe(string(domain.EventTypeNewConnection))
	h.eb.Subscribe(string(domain.EventTypeSendMessage))

	routes.RegisterBaseRoutes(h.engine)
	if err := h.engine.Run(":8080"); err != nil {
		h.logger.Warn("http server stopped", "error", err)
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			h.logger.Warn("error while serving http server", "error", err)
		}
	}()

	<-ctx.Done()
	return nil
}

// Stop shutdown gracefully the http service.
func (h *HTTP) Stop(ctx context.Context) error {
	return nil
}

// Name return the service name.
func (h *HTTP) Name() string {
	return "hopper-http-server"
}

// RegisterEventBus attach the event bus to the service.
func (h *HTTP) RegisterEventBus(eb domain.IEventBus) {
	h.eb = eb
}
