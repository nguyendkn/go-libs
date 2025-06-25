// Package websocket provides a comprehensive WebSocket library for Go
// with support for both client and server implementations, room management,
// authentication, rate limiting, and Gin framework integration.
package websocket

import (
	"net/http"
	"time"
)

// Version information
const (
	Version = "1.0.0"
	Author  = "Go-Libs WebSocket Team"
)

// Quick start functions for common use cases

// QuickServer creates a WebSocket server with sensible defaults
func QuickServer(addr string) Server {
	options := &ServerOptions{
		Addr:              addr,
		Path:              "/ws",
		EnableCompression: true,
		EnableMetrics:     true,
		MetricsPath:       "/metrics",
		PingInterval:      DefaultPingInterval,
		PongTimeout:       DefaultPongTimeout,
		MaxMessageSize:    DefaultMaxMessageSize,
		MessageQueueSize:  DefaultMessageQueueSize,
		RateLimit:         DefaultRateLimit,
		RateBurst:         DefaultRateBurst,
		ShutdownTimeout:   DefaultShutdownTimeout,
	}

	return NewServer(addr, options)
}

// QuickClient creates a WebSocket client with sensible defaults
func QuickClient(url string) Client {
	options := &ClientOptions{
		URL:                  url,
		AutoReconnect:        true,
		ReconnectInterval:    DefaultReconnectInterval,
		MaxReconnectAttempts: DefaultMaxReconnectAttempts,
		PingInterval:         DefaultPingInterval,
		PongTimeout:          DefaultPongTimeout,
		WriteTimeout:         DefaultWriteTimeout,
		ReadTimeout:          DefaultReadTimeout,
		MessageQueueSize:     DefaultMessageQueueSize,
		MaxMessageSize:       DefaultMaxMessageSize,
		RateLimit:            DefaultRateLimit,
		RateBurst:            DefaultRateBurst,
		EnableCompression:    false,
	}

	return NewClient(url, options)
}

// QuickGinServer creates a Gin-integrated WebSocket server
func QuickGinServer(addr string) *GinWebSocketServer {
	options := &ServerOptions{
		Addr:              addr,
		Path:              "/ws",
		EnableCompression: true,
		EnableMetrics:     true,
		MetricsPath:       "/metrics",
		AllowedOrigins:    []string{"*"},
		PingInterval:      DefaultPingInterval,
		PongTimeout:       DefaultPongTimeout,
		MaxMessageSize:    DefaultMaxMessageSize,
		RateLimit:         DefaultRateLimit,
		RateBurst:         DefaultRateBurst,
	}

	return NewGinWebSocketServer(addr, options)
}

// Builder patterns for advanced configuration

// ServerBuilder provides a fluent interface for building servers
type ServerBuilder struct {
	options *ServerOptions
}

// NewServerBuilder creates a new server builder
func NewServerBuilder(addr string) *ServerBuilder {
	return &ServerBuilder{
		options: &ServerOptions{
			Addr:              addr,
			Path:              "/ws",
			ReadTimeout:       DefaultReadTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			IdleTimeout:       DefaultIdleTimeout,
			MaxMessageSize:    DefaultMaxMessageSize,
			MessageQueueSize:  DefaultMessageQueueSize,
			PingInterval:      DefaultPingInterval,
			PongTimeout:       DefaultPongTimeout,
			RateLimit:         DefaultRateLimit,
			RateBurst:         DefaultRateBurst,
			RateWindow:        DefaultRateWindow,
			EnableCompression: false,
			CompressionLevel:  1,
			ShutdownTimeout:   DefaultShutdownTimeout,
			EnableMetrics:     true,
			MetricsPath:       "/metrics",
		},
	}
}

// WithPath sets the WebSocket path
func (b *ServerBuilder) WithPath(path string) *ServerBuilder {
	b.options.Path = path
	return b
}

// WithCompression enables compression
func (b *ServerBuilder) WithCompression(level int) *ServerBuilder {
	b.options.EnableCompression = true
	b.options.CompressionLevel = level
	return b
}

// WithAuth enables authentication
func (b *ServerBuilder) WithAuth(handler func(*http.Request) (*AuthInfo, error)) *ServerBuilder {
	b.options.AuthRequired = true
	b.options.AuthHandler = handler
	return b
}

// WithJWTAuth enables JWT authentication
func (b *ServerBuilder) WithJWTAuth(secret, issuer string, expiration time.Duration) *ServerBuilder {
	b.options.AuthRequired = true
	b.options.JWTSecret = secret
	b.options.AuthHandler = CreateJWTAuthHandler(secret, issuer, expiration)
	return b
}

// WithRateLimit sets rate limiting
func (b *ServerBuilder) WithRateLimit(limit, burst int) *ServerBuilder {
	b.options.RateLimit = limit
	b.options.RateBurst = burst
	return b
}

// WithCORS sets CORS configuration
func (b *ServerBuilder) WithCORS(origins, headers []string) *ServerBuilder {
	b.options.AllowedOrigins = origins
	b.options.AllowedHeaders = headers
	return b
}

// WithTLS enables TLS
func (b *ServerBuilder) WithTLS(certFile, keyFile string) *ServerBuilder {
	b.options.TLSCertFile = certFile
	b.options.TLSKeyFile = keyFile
	return b
}

// WithMetrics configures metrics
func (b *ServerBuilder) WithMetrics(enabled bool, path string) *ServerBuilder {
	b.options.EnableMetrics = enabled
	b.options.MetricsPath = path
	return b
}

// WithTimeouts sets various timeouts
func (b *ServerBuilder) WithTimeouts(read, write, idle, shutdown time.Duration) *ServerBuilder {
	b.options.ReadTimeout = read
	b.options.WriteTimeout = write
	b.options.IdleTimeout = idle
	b.options.ShutdownTimeout = shutdown
	return b
}

// WithHeartbeat sets heartbeat configuration
func (b *ServerBuilder) WithHeartbeat(pingInterval, pongTimeout time.Duration) *ServerBuilder {
	b.options.PingInterval = pingInterval
	b.options.PongTimeout = pongTimeout
	return b
}

// WithMessageLimits sets message size and queue limits
func (b *ServerBuilder) WithMessageLimits(maxSize int64, queueSize int) *ServerBuilder {
	b.options.MaxMessageSize = maxSize
	b.options.MessageQueueSize = queueSize
	return b
}

// Build creates the server
func (b *ServerBuilder) Build() Server {
	return NewServer(b.options.Addr, b.options)
}

// ClientBuilder provides a fluent interface for building clients
type ClientBuilder struct {
	options *ClientOptions
}

// NewClientBuilder creates a new client builder
func NewClientBuilder(url string) *ClientBuilder {
	return &ClientBuilder{
		options: &ClientOptions{
			URL:                  url,
			AutoReconnect:        true,
			ReconnectInterval:    DefaultReconnectInterval,
			MaxReconnectAttempts: DefaultMaxReconnectAttempts,
			ReconnectBackoff:     time.Second,
			PingInterval:         DefaultPingInterval,
			PongTimeout:          DefaultPongTimeout,
			WriteTimeout:         DefaultWriteTimeout,
			ReadTimeout:          DefaultReadTimeout,
			MessageQueueSize:     DefaultMessageQueueSize,
			MaxMessageSize:       DefaultMaxMessageSize,
			RateLimit:            DefaultRateLimit,
			RateBurst:            DefaultRateBurst,
			RateWindow:           DefaultRateWindow,
			EnableCompression:    false,
			Headers:              make(map[string]string),
		},
	}
}

// WithHeaders sets custom headers
func (b *ClientBuilder) WithHeaders(headers map[string]string) *ClientBuilder {
	b.options.Headers = headers
	return b
}

// WithAuth sets authentication
func (b *ClientBuilder) WithAuth(token, header string) *ClientBuilder {
	b.options.AuthToken = token
	b.options.AuthHeader = header
	return b
}

// WithAuthCallback sets authentication callback
func (b *ClientBuilder) WithAuthCallback(callback func() (string, error)) *ClientBuilder {
	b.options.AuthCallback = callback
	return b
}

// WithReconnect configures reconnection
func (b *ClientBuilder) WithReconnect(enabled bool, interval time.Duration, maxAttempts int) *ClientBuilder {
	b.options.AutoReconnect = enabled
	b.options.ReconnectInterval = interval
	b.options.MaxReconnectAttempts = maxAttempts
	return b
}

// WithCompression enables compression
func (b *ClientBuilder) WithCompression() *ClientBuilder {
	b.options.EnableCompression = true
	return b
}

// WithRateLimit sets rate limiting
func (b *ClientBuilder) WithRateLimit(limit, burst int) *ClientBuilder {
	b.options.RateLimit = limit
	b.options.RateBurst = burst
	return b
}

// WithTimeouts sets various timeouts
func (b *ClientBuilder) WithTimeouts(read, write time.Duration) *ClientBuilder {
	b.options.ReadTimeout = read
	b.options.WriteTimeout = write
	return b
}

// WithHeartbeat sets heartbeat configuration
func (b *ClientBuilder) WithHeartbeat(pingInterval, pongTimeout time.Duration) *ClientBuilder {
	b.options.PingInterval = pingInterval
	b.options.PongTimeout = pongTimeout
	return b
}

// WithProxy sets proxy configuration
func (b *ClientBuilder) WithProxy(proxyURL string) *ClientBuilder {
	b.options.ProxyURL = proxyURL
	return b
}

// Build creates the client
func (b *ClientBuilder) Build() Client {
	return NewClient(b.options.URL, b.options)
}

// Utility functions

// GetVersion returns the library version
func GetVersion() string {
	return Version
}

// GetAuthor returns the library author
func GetAuthor() string {
	return Author
}

// IsValidMessageType checks if a message type is valid
func IsValidMessageType(messageType MessageType) bool {
	switch messageType {
	case TextMessage, BinaryMessage, CloseMessage, PingMessage, PongMessage:
		return true
	default:
		return false
	}
}

// IsValidConnectionState checks if a connection state is valid
func IsValidConnectionState(state ConnectionState) bool {
	switch state {
	case StateDisconnected, StateConnecting, StateConnected, StateReconnecting, StateClosed:
		return true
	default:
		return false
	}
}
