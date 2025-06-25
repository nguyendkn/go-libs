package httpclient

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

// Client interface định nghĩa HTTP client
type Client interface {
	// Core HTTP methods
	Get(url string) RequestBuilder
	Post(url string) RequestBuilder
	Put(url string) RequestBuilder
	Patch(url string) RequestBuilder
	Delete(url string) RequestBuilder
	Head(url string) RequestBuilder
	Options(url string) RequestBuilder

	// Generic request method
	Request(method HTTPMethod, url string) RequestBuilder

	// Execute request
	Do(req *Request) (*Response, error)
	DoWithContext(ctx context.Context, req *Request) (*Response, error)

	// Configuration
	SetBaseURL(url string) Client
	SetUserAgent(userAgent string) Client
	SetTimeout(timeout time.Duration) Client
	SetHeaders(headers map[string]string) Client
	SetAuth(auth *AuthConfig) Client

	// Middleware
	Use(middleware Middleware) Client

	// Clone creates a copy of the client
	Clone() Client

	// Close closes the client and releases resources
	Close() error
}

// RequestBuilder interface cho fluent API
type RequestBuilder interface {
	// URL and path
	URL(url string) RequestBuilder
	Path(path string) RequestBuilder
	Pathf(format string, args ...interface{}) RequestBuilder

	// Headers
	Header(key, value string) RequestBuilder
	Headers(headers map[string]string) RequestBuilder
	ContentType(contentType ContentType) RequestBuilder
	Accept(accept string) RequestBuilder
	UserAgent(userAgent string) RequestBuilder

	// Query parameters
	Query(key, value string) RequestBuilder
	QueryParams(params map[string]string) RequestBuilder
	QueryStruct(v interface{}) RequestBuilder

	// Body
	Body(body interface{}) RequestBuilder
	BodyReader(reader io.Reader) RequestBuilder
	JSON(v interface{}) RequestBuilder
	XML(v interface{}) RequestBuilder
	Form(data map[string]string) RequestBuilder
	FormData(data map[string][]string) RequestBuilder
	File(fieldName, fileName string, reader io.Reader) RequestBuilder

	// Authentication
	Auth(auth *AuthConfig) RequestBuilder
	BasicAuth(username, password string) RequestBuilder
	BearerToken(token string) RequestBuilder
	APIKey(key, value string) RequestBuilder

	// Request options
	Timeout(timeout time.Duration) RequestBuilder
	Context(ctx context.Context) RequestBuilder
	FollowRedirects(follow bool) RequestBuilder
	MaxRedirects(max int) RequestBuilder

	// Retry
	Retry(policy *RetryPolicy) RequestBuilder
	RetryAttempts(attempts int) RequestBuilder
	RetryDelay(delay time.Duration) RequestBuilder

	// Cache
	Cache(ttl time.Duration) RequestBuilder
	CacheKey(key string) RequestBuilder
	NoCache() RequestBuilder

	// Metadata
	Metadata(key string, value interface{}) RequestBuilder

	// Execute
	Send() (*Response, error)
	SendWithContext(ctx context.Context) (*Response, error)

	// Response helpers
	Expect(statusCode int) (*Response, error)
	ExpectJSON(v interface{}) (*Response, error)
	ExpectXML(v interface{}) (*Response, error)
	ExpectBytes() ([]byte, error)
	ExpectString() (string, error)

	// Build request without sending
	Build() (*Request, error)
}

// Middleware interface cho request/response processing
type Middleware interface {
	Process(req *Request, next Handler) (*Response, error)
}

// Handler function type
type Handler func(*Request) (*Response, error)

// MiddlewareFunc adapter
type MiddlewareFunc func(*Request, Handler) (*Response, error)

func (f MiddlewareFunc) Process(req *Request, next Handler) (*Response, error) {
	return f(req, next)
}

// Cache interface cho caching system
type Cache interface {
	Get(key string) (*Response, bool)
	Set(key string, response *Response, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Size() int
	Stats() CacheStats
}

// CacheStats thống kê cache
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	Entries     int       `json:"entries"`
	Size        int64     `json:"size"`
	HitRate     float64   `json:"hitRate"`
	LastCleanup time.Time `json:"lastCleanup"`
}

// CircuitBreaker interface cho circuit breaker pattern
type CircuitBreaker interface {
	Execute(req *Request, fn func(*Request) (*Response, error)) (*Response, error)
	State() CircuitState
	Metrics() CircuitMetrics
	Reset()
}

// CircuitState trạng thái circuit breaker
type CircuitState int

const (
	CircuitStateClosed CircuitState = iota
	CircuitStateOpen
	CircuitStateHalfOpen
)

// CircuitMetrics metrics của circuit breaker
type CircuitMetrics struct {
	Requests     int64        `json:"requests"`
	Successes    int64        `json:"successes"`
	Failures     int64        `json:"failures"`
	Timeouts     int64        `json:"timeouts"`
	LastFailure  time.Time    `json:"lastFailure"`
	LastSuccess  time.Time    `json:"lastSuccess"`
	FailureRate  float64      `json:"failureRate"`
	State        CircuitState `json:"state"`
	StateChanged time.Time    `json:"stateChanged"`
}

// RateLimiter interface cho rate limiting
type RateLimiter interface {
	Allow(key string) bool
	Wait(key string) time.Duration
	Reserve(key string) Reservation
	Limit() float64
	Burst() int
}

// Reservation đại diện cho rate limit reservation
type Reservation interface {
	OK() bool
	Delay() time.Duration
	Cancel()
}

// Authenticator interface cho authentication
type Authenticator interface {
	Authenticate(req *Request) error
	Type() AuthType
	IsExpired() bool
	Refresh() error
}

// Logger interface cho logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

// Field đại diện cho log field
type Field struct {
	Key   string
	Value interface{}
}

// Metrics interface cho metrics collection
type Metrics interface {
	RecordRequest(req *Request)
	RecordResponse(req *Request, resp *Response, duration time.Duration)
	RecordError(req *Request, err error, duration time.Duration)
	GetStats() MetricsStats
}

// MetricsStats thống kê metrics
type MetricsStats struct {
	TotalRequests  int64         `json:"totalRequests"`
	TotalResponses int64         `json:"totalResponses"`
	TotalErrors    int64         `json:"totalErrors"`
	AverageLatency time.Duration `json:"averageLatency"`
	P50Latency     time.Duration `json:"p50Latency"`
	P95Latency     time.Duration `json:"p95Latency"`
	P99Latency     time.Duration `json:"p99Latency"`
	RequestsPerSec float64       `json:"requestsPerSec"`
	ErrorRate      float64       `json:"errorRate"`
	StatusCodes    map[int]int64 `json:"statusCodes"`
	LastUpdated    time.Time     `json:"lastUpdated"`
}

// Tracer interface cho distributed tracing
type Tracer interface {
	StartSpan(req *Request) Span
	InjectHeaders(span Span, headers map[string]string)
	ExtractSpan(headers map[string]string) SpanContext
}

// Span đại diện cho tracing span
type Span interface {
	SetTag(key string, value interface{})
	SetError(err error)
	Finish()
	Context() SpanContext
}

// SpanContext đại diện cho span context
type SpanContext interface {
	TraceID() string
	SpanID() string
	IsSampled() bool
}

// HealthChecker interface cho health checking
type HealthChecker interface {
	Check(ctx context.Context) HealthStatus
	Name() string
}

// HealthStatus trạng thái health check
type HealthStatus struct {
	Status    string                 `json:"status"` // "healthy", "unhealthy", "degraded"
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
}

// Transport interface cho custom transport
type Transport interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

// Interceptor interface cho request/response interception
type Interceptor interface {
	BeforeRequest(req *Request) error
	AfterResponse(req *Request, resp *Response) error
	OnError(req *Request, err error) error
}

// Serializer interface cho body serialization
type Serializer interface {
	Serialize(v interface{}) ([]byte, error)
	Deserialize(data []byte, v interface{}) error
	ContentType() string
}

// Compressor interface cho response compression
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
	Encoding() string
}

// RetryDecider interface cho custom retry logic
type RetryDecider interface {
	ShouldRetry(req *Request, resp *Response, err error, attempt int) bool
	NextDelay(req *Request, resp *Response, err error, attempt int) time.Duration
}

// LoadBalancer interface cho load balancing
type LoadBalancer interface {
	Next(req *Request) (string, error)
	MarkUnhealthy(endpoint string)
	MarkHealthy(endpoint string)
	Endpoints() []string
}

// ConnectionPool interface cho connection pooling
type ConnectionPool interface {
	Get(network, addr string) (net.Conn, error)
	Put(conn net.Conn) error
	Close() error
	Stats() PoolStats
}

// PoolStats thống kê connection pool
type PoolStats struct {
	Active    int `json:"active"`
	Idle      int `json:"idle"`
	Total     int `json:"total"`
	MaxActive int `json:"maxActive"`
	MaxIdle   int `json:"maxIdle"`
}

// ResponseProcessor interface cho response processing
type ResponseProcessor interface {
	Process(resp *Response) error
	ContentTypes() []string
}

// RequestValidator interface cho request validation
type RequestValidator interface {
	Validate(req *Request) error
}

// ErrorHandler interface cho error handling
type ErrorHandler interface {
	Handle(req *Request, err error) (*Response, error)
	ShouldRetry(req *Request, err error) bool
}

// ConfigProvider interface cho dynamic configuration
type ConfigProvider interface {
	GetConfig() *ClientConfig
	Watch(callback func(*ClientConfig))
	Close() error
}
