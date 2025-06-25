package websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType định nghĩa các loại message
type MessageType int

const (
	TextMessage   MessageType = 1
	BinaryMessage MessageType = 2
	CloseMessage  MessageType = 8
	PingMessage   MessageType = 9
	PongMessage   MessageType = 10
)

// ConnectionState định nghĩa trạng thái kết nối
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
	StateClosed
)

// Message đại diện cho một message WebSocket
type Message struct {
	Type      MessageType `json:"type"`
	Data      []byte      `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ClientID  string      `json:"client_id,omitempty"`
	Room      string      `json:"room,omitempty"`
}

// ClientInfo chứa thông tin về client
type ClientInfo struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id,omitempty"`
	IP          string            `json:"ip"`
	UserAgent   string            `json:"user_agent"`
	Headers     map[string]string `json:"headers,omitempty"`
	Rooms       []string          `json:"rooms,omitempty"`
	ConnectedAt time.Time         `json:"connected_at"`
	LastPing    time.Time         `json:"last_ping"`
}

// ConnectionMetrics chứa metrics của kết nối
type ConnectionMetrics struct {
	MessagesSent     int64     `json:"messages_sent"`
	MessagesReceived int64     `json:"messages_received"`
	BytesSent        int64     `json:"bytes_sent"`
	BytesReceived    int64     `json:"bytes_received"`
	LastActivity     time.Time `json:"last_activity"`
	Errors           int64     `json:"errors"`
}

// ClientOptions cấu hình cho WebSocket client
type ClientOptions struct {
	// Connection settings
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Subprotocols []string          `json:"subprotocols,omitempty"`

	// Reconnection settings
	AutoReconnect        bool          `json:"auto_reconnect"`
	ReconnectInterval    time.Duration `json:"reconnect_interval"`
	MaxReconnectAttempts int           `json:"max_reconnect_attempts"`
	ReconnectBackoff     time.Duration `json:"reconnect_backoff"`

	// Heartbeat settings
	PingInterval time.Duration `json:"ping_interval"`
	PongTimeout  time.Duration `json:"pong_timeout"`

	// Message settings
	WriteTimeout     time.Duration `json:"write_timeout"`
	ReadTimeout      time.Duration `json:"read_timeout"`
	MessageQueueSize int           `json:"message_queue_size"`
	MaxMessageSize   int64         `json:"max_message_size"`

	// Compression
	EnableCompression bool `json:"enable_compression"`

	// Authentication
	AuthToken    string                 `json:"auth_token,omitempty"`
	AuthHeader   string                 `json:"auth_header,omitempty"`
	AuthCallback func() (string, error) `json:"-"`

	// Rate limiting
	RateLimit  int           `json:"rate_limit"`  // messages per second
	RateBurst  int           `json:"rate_burst"`  // burst size
	RateWindow time.Duration `json:"rate_window"` // rate window

	// TLS settings
	TLSConfig interface{} `json:"-"`

	// Proxy settings
	ProxyURL string `json:"proxy_url,omitempty"`

	// Custom dialer
	Dialer *websocket.Dialer `json:"-"`
}

// ServerOptions cấu hình cho WebSocket server
type ServerOptions struct {
	// Server settings
	Addr         string        `json:"addr"`
	Path         string        `json:"path"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`

	// WebSocket settings
	CheckOrigin       func(*http.Request) bool `json:"-"`
	Subprotocols      []string                 `json:"subprotocols,omitempty"`
	EnableCompression bool                     `json:"enable_compression"`
	CompressionLevel  int                      `json:"compression_level"`

	// Message settings
	MaxMessageSize   int64 `json:"max_message_size"`
	MessageQueueSize int   `json:"message_queue_size"`

	// Heartbeat settings
	PingInterval time.Duration `json:"ping_interval"`
	PongTimeout  time.Duration `json:"pong_timeout"`

	// Rate limiting
	RateLimit  int           `json:"rate_limit"`  // messages per second per client
	RateBurst  int           `json:"rate_burst"`  // burst size per client
	RateWindow time.Duration `json:"rate_window"` // rate window

	// Authentication
	AuthRequired bool                                   `json:"auth_required"`
	AuthHandler  func(*http.Request) (*AuthInfo, error) `json:"-"`
	JWTSecret    string                                 `json:"jwt_secret,omitempty"`

	// CORS settings
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	AllowedHeaders []string `json:"allowed_headers,omitempty"`

	// TLS settings
	TLSCertFile string `json:"tls_cert_file,omitempty"`
	TLSKeyFile  string `json:"tls_key_file,omitempty"`

	// Monitoring
	EnableMetrics bool   `json:"enable_metrics"`
	MetricsPath   string `json:"metrics_path,omitempty"`

	// Graceful shutdown
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

// AuthInfo chứa thông tin authentication
type AuthInfo struct {
	UserID   string                 `json:"user_id"`
	Username string                 `json:"username,omitempty"`
	Roles    []string               `json:"roles,omitempty"`
	Claims   map[string]interface{} `json:"claims,omitempty"`
	Token    string                 `json:"token,omitempty"`
}

// RoomOptions cấu hình cho room
type RoomOptions struct {
	MaxClients     int           `json:"max_clients"`
	RequireAuth    bool          `json:"require_auth"`
	AllowedRoles   []string      `json:"allowed_roles,omitempty"`
	MessageHistory int           `json:"message_history"` // số message lưu lại
	TTL            time.Duration `json:"ttl"`             // thời gian sống của room
}

// EventType định nghĩa các loại event
type EventType string

const (
	EventConnect    EventType = "connect"
	EventDisconnect EventType = "disconnect"
	EventMessage    EventType = "message"
	EventError      EventType = "error"
	EventJoinRoom   EventType = "join_room"
	EventLeaveRoom  EventType = "leave_room"
	EventPing       EventType = "ping"
	EventPong       EventType = "pong"
)

// Event đại diện cho một event trong hệ thống
type Event struct {
	Type      EventType   `json:"type"`
	ClientID  string      `json:"client_id"`
	Room      string      `json:"room,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Error     error       `json:"error,omitempty"`
}

// Default values
const (
	DefaultPingInterval         = 30 * time.Second
	DefaultPongTimeout          = 10 * time.Second
	DefaultReconnectInterval    = 5 * time.Second
	DefaultMaxReconnectAttempts = 10
	DefaultWriteTimeout         = 10 * time.Second
	DefaultReadTimeout          = 60 * time.Second
	DefaultMessageQueueSize     = 1000
	DefaultMaxMessageSize       = 1024 * 1024 // 1MB
	DefaultRateLimit            = 100         // messages per second
	DefaultRateBurst            = 200
	DefaultRateWindow           = time.Second
	DefaultShutdownTimeout      = 30 * time.Second
	DefaultIdleTimeout          = 60 * time.Second
)

// Error definitions
var (
	ErrConnectionClosed     = websocket.ErrCloseSent
	ErrMessageTooLarge      = &websocket.CloseError{Code: websocket.CloseMessageTooBig, Text: "message too large"}
	ErrRateLimitExceeded    = &websocket.CloseError{Code: websocket.ClosePolicyViolation, Text: "rate limit exceeded"}
	ErrAuthenticationFailed = &websocket.CloseError{Code: websocket.CloseUnsupportedData, Text: "authentication failed"}
	ErrUnauthorized         = &websocket.CloseError{Code: websocket.CloseUnsupportedData, Text: "unauthorized"}
	ErrRoomNotFound         = &websocket.CloseError{Code: websocket.CloseUnsupportedData, Text: "room not found"}
	ErrRoomFull             = &websocket.CloseError{Code: websocket.ClosePolicyViolation, Text: "room is full"}
)
