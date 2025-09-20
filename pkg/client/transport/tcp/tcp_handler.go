package tcp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type HealthStatus int32

const (
	StatusDisconnected HealthStatus = iota
	StatusConnecting
	StatusConnected
	StatusDegraded
	StatusFailed
)

type CircuitState int32

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// TCPClient provides a production-grade TCP client service with health monitoring and auto-reconnection
type TCPClient struct {
	conn   net.Conn
	connMu sync.RWMutex

	logger *slog.Logger
	config struct {
		port         int
		address      string
		maxRetries   int
		maxRetryWait time.Duration
		retryDelay   time.Duration
		retryBackoff float64

		healthCheckInterval time.Duration
		circuitFailures     int
		circuitTimeout      time.Duration
		reconnectInterval   time.Duration
	}

	health struct {
		status           int32
		lastConnected    time.Time
		lastError        error
		consecutiveErrors int32
		totalConnections int64
		totalErrors      int64
	}

	circuit struct {
		state       int32
		failures    int32
		lastFailure time.Time
		nextRetry   time.Time
	}

	cancel  context.CancelFunc
	done    chan struct{}
	errChan chan error

	wg sync.WaitGroup
}

type Option func(*TCPClient)

// WithLogger sets the logger for the TCP client
func WithLogger(logger *slog.Logger) Option {
	return func(c *TCPClient) {
		c.logger = logger
	}
}

// WithAddress sets the target address and port
func WithAddress(address string, port int) Option {
	return func(c *TCPClient) {
		c.config.address = address
		c.config.port = port
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(maxRetries int, delay time.Duration, backoff float64, maxWait time.Duration) Option {
	return func(c *TCPClient) {
		c.config.maxRetries = maxRetries
		c.config.retryDelay = delay
		c.config.retryBackoff = backoff
		c.config.maxRetryWait = maxWait
	}
}

// WithHealthConfig sets health monitoring configuration
func WithHealthConfig(interval time.Duration, circuitFailures int, circuitTimeout time.Duration) Option {
	return func(c *TCPClient) {
		c.config.healthCheckInterval = interval
		c.config.circuitFailures = circuitFailures
		c.config.circuitTimeout = circuitTimeout
	}
}

// Start initializes and starts the TCP client service
func (t *TCPClient) Start(ctx context.Context) error {
	ctx, t.cancel = context.WithCancel(ctx)

	t.setHealthStatus(StatusConnecting)

	t.wg.Add(3)
	go t.connectionManager(ctx)
	go t.healthMonitor(ctx)
	go t.messageHandler(ctx)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-t.errChan:
		return err
	case <-time.After(10 * time.Second):
		if t.IsHealthy() {
			return nil
		}
		return fmt.Errorf("failed to establish initial connection")
	}
}

// Stop gracefully shuts down the TCP client service
func (t *TCPClient) Stop(ctx context.Context) error {
	if t.cancel != nil {
		t.cancel()
	}

	done := make(chan struct{})
	go func() {
		t.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.closeConnection()
		close(t.done)
		return nil
	case <-ctx.Done():
		t.closeConnection()
		return ctx.Err()
	}
}

// Name returns the service name
func (t *TCPClient) Name() string {
	return "tcp_transport_service"
}

// IsHealthy returns true if the client is in a healthy state
func (t *TCPClient) IsHealthy() bool {
	status := HealthStatus(atomic.LoadInt32(&t.health.status))
	return status == StatusConnected
}

// IsDegraded returns true if the client is degraded but operational
func (t *TCPClient) IsDegraded() bool {
	status := HealthStatus(atomic.LoadInt32(&t.health.status))
	return status == StatusDegraded
}

// GetHealthStatus returns the current health status
func (t *TCPClient) GetHealthStatus() HealthStatus {
	return HealthStatus(atomic.LoadInt32(&t.health.status))
}

// GetHealthMetrics returns health metrics
func (t *TCPClient) GetHealthMetrics() map[string]interface{} {
	t.connMu.RLock()
	defer t.connMu.RUnlock()

	return map[string]interface{}{
		"status":             t.GetHealthStatus(),
		"last_connected":     t.health.lastConnected,
		"consecutive_errors": atomic.LoadInt32(&t.health.consecutiveErrors),
		"total_connections":  atomic.LoadInt64(&t.health.totalConnections),
		"total_errors":       atomic.LoadInt64(&t.health.totalErrors),
		"circuit_state":      CircuitState(atomic.LoadInt32(&t.circuit.state)),
		"circuit_failures":   atomic.LoadInt32(&t.circuit.failures),
	}
}

func (t *TCPClient) connectionManager(ctx context.Context) {
	defer t.wg.Done()

	ticker := time.NewTicker(t.config.reconnectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !t.isConnected() && t.shouldAttemptConnection() {
				t.attemptConnection(ctx)
			}
		}
	}
}

func (t *TCPClient) healthMonitor(ctx context.Context) {
	defer t.wg.Done()

	ticker := time.NewTicker(t.config.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.performHealthCheck()
		}
	}
}

func (t *TCPClient) messageHandler(ctx context.Context) {
	defer t.wg.Done()

	buffer := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn := t.getConnection()
		if conn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.Read(buffer)

		if err != nil {
			t.handleConnectionError(err)
			continue
		}

		if n > 0 {
			t.resetConsecutiveErrors()
			t.logger.Info("received", "bytes", n, "data", string(buffer[:n]))
			continue
		}
	}
}

func (t *TCPClient) attemptConnection(ctx context.Context) {
	if !t.shouldAttemptConnection() {
		return
	}

	t.setHealthStatus(StatusConnecting)

	addr := net.JoinHostPort(t.config.address, strconv.Itoa(t.config.port))

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		t.handleConnectionError(err)
		return
	}

	t.setConnection(conn)
	t.setHealthStatus(StatusConnected)
	t.resetCircuit()

	t.health.lastConnected = time.Now()
	atomic.AddInt64(&t.health.totalConnections, 1)

	t.logger.Info("connected successfully", "address", addr)
}

func (t *TCPClient) handleConnectionError(err error) {
	t.closeConnection()

	consecutive := atomic.AddInt32(&t.health.consecutiveErrors, 1)
	atomic.AddInt64(&t.health.totalErrors, 1)

	t.health.lastError = err
	t.recordCircuitFailure()

	if consecutive < 3 {
		t.setHealthStatus(StatusDegraded)
	} else {
		t.setHealthStatus(StatusFailed)
	}

	t.logger.Warn("connection error",
		"error", err,
		"consecutive_errors", consecutive,
		"circuit_state", CircuitState(atomic.LoadInt32(&t.circuit.state)))
}

func (t *TCPClient) performHealthCheck() {
	conn := t.getConnection()
	if conn == nil {
		if t.GetHealthStatus() == StatusConnected {
			t.setHealthStatus(StatusDisconnected)
		}
		return
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
}

func (t *TCPClient) shouldAttemptConnection() bool {
	circuitState := CircuitState(atomic.LoadInt32(&t.circuit.state))

	switch circuitState {
	case CircuitOpen:
		return time.Now().After(t.circuit.nextRetry)
	case CircuitHalfOpen:
		return true
	default:
		return true
	}
}

func (t *TCPClient) recordCircuitFailure() {
	failures := atomic.AddInt32(&t.circuit.failures, 1)
	t.circuit.lastFailure = time.Now()

	if failures >= int32(t.config.circuitFailures) {
		atomic.StoreInt32(&t.circuit.state, int32(CircuitOpen))
		t.circuit.nextRetry = time.Now().Add(t.config.circuitTimeout)
		t.logger.Warn("circuit breaker opened", "failures", failures)
	}
}

func (t *TCPClient) resetCircuit() {
	atomic.StoreInt32(&t.circuit.state, int32(CircuitClosed))
	atomic.StoreInt32(&t.circuit.failures, 0)
}

func (t *TCPClient) resetConsecutiveErrors() {
	atomic.StoreInt32(&t.health.consecutiveErrors, 0)
}

func (t *TCPClient) setHealthStatus(status HealthStatus) {
	atomic.StoreInt32(&t.health.status, int32(status))
}

func (t *TCPClient) isConnected() bool {
	t.connMu.RLock()
	defer t.connMu.RUnlock()
	return t.conn != nil
}

func (t *TCPClient) getConnection() net.Conn {
	t.connMu.RLock()
	defer t.connMu.RUnlock()
	return t.conn
}

func (t *TCPClient) setConnection(conn net.Conn) {
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if t.conn != nil {
		t.conn.Close()
	}
	t.conn = conn
}

func (t *TCPClient) closeConnection() {
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	t.setHealthStatus(StatusDisconnected)
}

// NewTCPClient creates a new production-grade TCP client service
func NewTCPClient(opts ...Option) *TCPClient {
	t := &TCPClient{
		done:    make(chan struct{}),
		errChan: make(chan error, 1),
		config: struct {
			port                int
			address             string
			maxRetries          int
			maxRetryWait        time.Duration
			retryDelay          time.Duration
			retryBackoff        float64
			healthCheckInterval time.Duration
			circuitFailures     int
			circuitTimeout      time.Duration
			reconnectInterval   time.Duration
		}{
			port:                5672,
			address:             "127.0.0.1",
			maxRetries:          10,
			maxRetryWait:        300 * time.Second,
			retryDelay:          1 * time.Second,
			retryBackoff:        1.5,
			healthCheckInterval: 30 * time.Second,
			circuitFailures:     5,
			circuitTimeout:      60 * time.Second,
			reconnectInterval:   5 * time.Second,
		},
	}

	atomic.StoreInt32(&t.health.status, int32(StatusDisconnected))
	atomic.StoreInt32(&t.circuit.state, int32(CircuitClosed))

	for _, opt := range opts {
		opt(t)
	}

	if t.logger == nil {
		t.logger = slog.Default()
	}

	return t
}
