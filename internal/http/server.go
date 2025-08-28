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

func WithContext(ctx context.CancelFunc) Option {
	return func(h *HTTP) {
		h.cancel = ctx
	}
}

func WithEventBus(eb domain.IEventBus) Option {
	return func(h *HTTP) {
		h.eb = eb
	}
}

func WithConfiguration(cfg *config.Configuration) Option {
	return func(h *HTTP) {
		h.config = cfg
	}
}

func WithEngine(engine *gin.Engine) Option {
	return func(h *HTTP) {
		h.engine = engine
	}
}

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

func (h *HTTP) Run(ctx context.Context) error {
	h.logger.Info("starting the http server component", "name", h.Name())
	h.eb.Subscribe(string(domain.EventTypeNewConnection))
	h.eb.Subscribe(string(domain.EventTypeSendMessage))

	routes.RegisterBaseRoutes(h.engine)
	h.engine.Run(":8080")

	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			h.logger.Warn("error while servint http server", "error", err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (h *HTTP) Stop(ctx context.Context) error {
	return nil
}

func (h *HTTP) Name() string {
	return "hopper-http-server"
}
