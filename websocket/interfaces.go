package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client interface định nghĩa các phương thức cho WebSocket client
type Client interface {
	// Connection management
	Connect() error
	ConnectWithContext(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
	GetState() ConnectionState

	// Message handling
	Send(data []byte) error
	SendText(text string) error
	SendJSON(v interface{}) error
	SendWithType(messageType MessageType, data []byte) error

	// Event handlers
	OnConnect(handler func())
	OnDisconnect(handler func(error))
	OnMessage(handler func([]byte))
	OnTextMessage(handler func(string))
	OnBinaryMessage(handler func([]byte))
	OnError(handler func(error))
	OnReconnect(handler func(attempt int))

	// Client info
	ID() string
	Info() *ClientInfo
	Metrics() *ConnectionMetrics

	// Configuration
	SetOptions(options *ClientOptions)
	GetOptions() *ClientOptions

	// Lifecycle
	Close() error
}

// Server interface định nghĩa các phương thức cho WebSocket server
type Server interface {
	// Server lifecycle
	Start() error
	StartWithContext(ctx context.Context) error
	Stop() error
	Shutdown(ctx context.Context) error

	// Client management
	GetClient(id string) (ServerClient, bool)
	GetClients() []ServerClient
	GetClientCount() int
	DisconnectClient(id string) error

	// Room management
	CreateRoom(name string, options *RoomOptions) error
	DeleteRoom(name string) error
	GetRoom(name string) (Room, bool)
	GetRooms() []Room

	// Broadcasting
	Broadcast(data []byte) error
	BroadcastText(text string) error
	BroadcastJSON(v interface{}) error
	BroadcastToRoom(room string, data []byte) error
	BroadcastToRoomText(room string, text string) error
	BroadcastToRoomJSON(room string, v interface{}) error

	// Event handlers
	OnConnect(handler func(ServerClient))
	OnDisconnect(handler func(ServerClient, error))
	OnMessage(handler func(ServerClient, []byte))
	OnTextMessage(handler func(ServerClient, string))
	OnBinaryMessage(handler func(ServerClient, []byte))
	OnError(handler func(ServerClient, error))
	OnRoomJoin(handler func(ServerClient, string))
	OnRoomLeave(handler func(ServerClient, string))

	// Configuration
	SetOptions(options *ServerOptions)
	GetOptions() *ServerOptions

	// Monitoring
	GetMetrics() *ServerMetrics
	GetHealth() *HealthStatus

	// HTTP handler
	GetHTTPHandler() http.Handler
	GetUpgrader() *websocket.Upgrader
}

// ServerClient interface định nghĩa client từ phía server
type ServerClient interface {
	// Basic info
	ID() string
	Info() *ClientInfo
	Metrics() *ConnectionMetrics

	// Connection state
	IsConnected() bool
	GetState() ConnectionState

	// Message sending
	Send(data []byte) error
	SendText(text string) error
	SendJSON(v interface{}) error
	SendWithType(messageType MessageType, data []byte) error

	// Room management
	JoinRoom(room string) error
	LeaveRoom(room string) error
	GetRooms() []string
	IsInRoom(room string) bool

	// Authentication
	GetAuth() *AuthInfo
	SetAuth(auth *AuthInfo)
	IsAuthenticated() bool

	// Rate limiting
	IsRateLimited() bool
	GetRateLimit() (int, int) // current, limit

	// Connection management
	Ping() error
	Close() error
	CloseWithCode(code int, text string) error

	// Context
	GetContext() context.Context
	SetContext(ctx context.Context)

	// Custom data
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
}

// Room interface định nghĩa room/channel
type Room interface {
	// Basic info
	Name() string
	Options() *RoomOptions

	// Client management
	AddClient(client ServerClient) error
	RemoveClient(clientID string) error
	GetClient(clientID string) (ServerClient, bool)
	GetClients() []ServerClient
	GetClientCount() int
	HasClient(clientID string) bool

	// Broadcasting
	Broadcast(data []byte) error
	BroadcastText(text string) error
	BroadcastJSON(v interface{}) error
	BroadcastExcept(data []byte, excludeClientID string) error

	// Message history
	GetMessageHistory() []*Message
	AddToHistory(msg *Message)
	ClearHistory()

	// Room state
	IsEmpty() bool
	IsFull() bool

	// Lifecycle
	Close() error

	// Events
	OnClientJoin(handler func(ServerClient))
	OnClientLeave(handler func(ServerClient))
	OnMessage(handler func(ServerClient, []byte))
	OnEmpty(handler func())
}

// Hub interface định nghĩa message hub
type Hub interface {
	// Client management
	RegisterClient(client ServerClient)
	UnregisterClient(clientID string)
	GetClient(clientID string) (ServerClient, bool)
	GetClients() []ServerClient
	GetClientCount() int

	// Room management
	CreateRoom(name string, options *RoomOptions) error
	DeleteRoom(name string) error
	GetRoom(name string) (Room, bool)
	GetRooms() []Room
	JoinRoom(clientID, roomName string) error
	LeaveRoom(clientID, roomName string) error

	// Broadcasting
	Broadcast(msg *Message) error
	BroadcastToRoom(roomName string, msg *Message) error
	BroadcastToClient(clientID string, msg *Message) error

	// Event handling
	HandleEvent(event *Event)

	// Lifecycle
	Start() error
	Stop() error

	// Monitoring
	GetMetrics() *HubMetrics
}

// RateLimiter interface định nghĩa rate limiter
type RateLimiter interface {
	Allow(clientID string) bool
	AllowN(clientID string, n int) bool
	Reset(clientID string)
	GetLimit(clientID string) (current int, limit int)
	SetLimit(clientID string, limit int, burst int)
	Cleanup() // cleanup expired entries
}

// Authenticator interface định nghĩa authentication
type Authenticator interface {
	Authenticate(req *http.Request) (*AuthInfo, error)
	ValidateToken(token string) (*AuthInfo, error)
	RefreshToken(token string) (string, error)
	RevokeToken(token string) error
}

// MessageQueue interface định nghĩa message queue
type MessageQueue interface {
	Enqueue(msg *Message) error
	Dequeue() (*Message, error)
	Peek() (*Message, error)
	Size() int
	IsEmpty() bool
	IsFull() bool
	Clear()
	Close() error
}

// EventBus interface định nghĩa event bus
type EventBus interface {
	Subscribe(eventType EventType, handler func(*Event))
	Unsubscribe(eventType EventType, handler func(*Event))
	Publish(event *Event)
	PublishAsync(event *Event)
	Close() error
}

// Logger interface định nghĩa logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// Metrics interfaces
type ServerMetrics struct {
	TotalConnections     int64     `json:"total_connections"`
	ActiveConnections    int64     `json:"active_connections"`
	TotalMessages        int64     `json:"total_messages"`
	TotalBytes           int64     `json:"total_bytes"`
	ErrorCount           int64     `json:"error_count"`
	RoomsCount           int       `json:"rooms_count"`
	AverageLatency       float64   `json:"average_latency_ms"`
	MessagesPerSecond    float64   `json:"messages_per_second"`
	ConnectionsPerSecond float64   `json:"connections_per_second"`
	Uptime               int64     `json:"uptime_seconds"`
	LastUpdated          time.Time `json:"last_updated"`
}

type HubMetrics struct {
	RegisteredClients int64                   `json:"registered_clients"`
	ActiveRooms       int                     `json:"active_rooms"`
	MessagesProcessed int64                   `json:"messages_processed"`
	EventsProcessed   int64                   `json:"events_processed"`
	QueueSize         int                     `json:"queue_size"`
	RoomMetrics       map[string]*RoomMetrics `json:"room_metrics"`
}

type RoomMetrics struct {
	Name         string    `json:"name"`
	ClientCount  int       `json:"client_count"`
	MessageCount int64     `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}

type HealthStatus struct {
	Status      string            `json:"status"` // "healthy", "degraded", "unhealthy"
	Checks      map[string]string `json:"checks"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      int64             `json:"uptime_seconds"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
}
