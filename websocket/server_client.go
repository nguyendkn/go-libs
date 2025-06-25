package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nguyendkn/go-libs/uuid"
	"golang.org/x/time/rate"
)

// serverClient implements the ServerClient interface
type serverClient struct {
	id      string
	conn    *websocket.Conn
	info    *ClientInfo
	metrics *ConnectionMetrics
	auth    *AuthInfo
	options *ServerOptions
	server  *wsServer // Reference to server for callbacks

	// State
	state int32 // atomic, ConnectionState

	// Context and data
	ctx    context.Context
	cancel context.CancelFunc
	data   map[string]any
	dataMu sync.RWMutex

	// Channels
	sendCh  chan *Message
	closeCh chan struct{}

	// Rate limiter
	rateLimiter *rate.Limiter

	// Rooms
	rooms   map[string]bool
	roomsMu sync.RWMutex

	// Heartbeat
	lastPong time.Time
	pongMu   sync.RWMutex

	// Wait group
	wg sync.WaitGroup
}

// newServerClient tạo một server client mới
func newServerClient(conn *websocket.Conn, req *http.Request, auth *AuthInfo, options *ServerOptions, server *wsServer) *serverClient {
	ctx, cancel := context.WithCancel(context.Background())

	client := &serverClient{
		id:      uuid.New().String(),
		conn:    conn,
		auth:    auth,
		options: options,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
		data:    make(map[string]any),
		sendCh:  make(chan *Message, options.MessageQueueSize),
		closeCh: make(chan struct{}),
		rooms:   make(map[string]bool),
		info: &ClientInfo{
			ID:          uuid.New().String(),
			IP:          getClientIP(req),
			UserAgent:   req.UserAgent(),
			Headers:     getHeaders(req),
			ConnectedAt: time.Now(),
			LastPing:    time.Now(),
		},
		metrics: &ConnectionMetrics{
			LastActivity: time.Now(),
		},
		rateLimiter: rate.NewLimiter(rate.Limit(options.RateLimit), options.RateBurst),
	}

	// Set user ID from auth if available
	if auth != nil {
		client.info.UserID = auth.UserID
	}

	client.pongMu.Lock()
	client.lastPong = time.Now()
	client.pongMu.Unlock()

	atomic.StoreInt32(&client.state, int32(StateConnected))

	return client
}

// Basic info methods
func (c *serverClient) ID() string {
	return c.id
}

func (c *serverClient) Info() *ClientInfo {
	return c.info
}

func (c *serverClient) Metrics() *ConnectionMetrics {
	return c.metrics
}

// Connection state methods
func (c *serverClient) IsConnected() bool {
	return atomic.LoadInt32(&c.state) == int32(StateConnected)
}

func (c *serverClient) GetState() ConnectionState {
	return ConnectionState(atomic.LoadInt32(&c.state))
}

// Message sending methods
func (c *serverClient) Send(data []byte) error {
	return c.SendWithType(BinaryMessage, data)
}

func (c *serverClient) SendText(text string) error {
	return c.SendWithType(TextMessage, []byte(text))
}

func (c *serverClient) SendJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return c.SendText(string(data))
}

func (c *serverClient) SendWithType(messageType MessageType, data []byte) error {
	if !c.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Check rate limit
	if !c.rateLimiter.Allow() {
		return ErrRateLimitExceeded
	}

	// Check message size
	if c.options.MaxMessageSize > 0 && int64(len(data)) > c.options.MaxMessageSize {
		return ErrMessageTooLarge
	}

	msg := &Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		ClientID:  c.id,
	}

	select {
	case c.sendCh <- msg:
		return nil
	case <-c.closeCh:
		return ErrConnectionClosed
	default:
		return fmt.Errorf("message queue is full")
	}
}

// Room management methods
func (c *serverClient) JoinRoom(room string) error {
	c.roomsMu.Lock()
	c.rooms[room] = true
	rooms := make([]string, 0, len(c.rooms))
	for r := range c.rooms {
		rooms = append(rooms, r)
	}
	c.info.Rooms = rooms
	c.roomsMu.Unlock()

	// Notify server about room join
	if c.server != nil {
		c.server.handleRoomJoin(c, room)
	}

	return nil
}

func (c *serverClient) LeaveRoom(room string) error {
	c.roomsMu.Lock()
	delete(c.rooms, room)
	rooms := make([]string, 0, len(c.rooms))
	for r := range c.rooms {
		rooms = append(rooms, r)
	}
	c.info.Rooms = rooms
	c.roomsMu.Unlock()

	// Notify server about room leave
	if c.server != nil {
		c.server.handleRoomLeave(c, room)
	}

	return nil
}

func (c *serverClient) GetRooms() []string {
	c.roomsMu.RLock()
	rooms := make([]string, 0, len(c.rooms))
	for room := range c.rooms {
		rooms = append(rooms, room)
	}
	c.roomsMu.RUnlock()
	return rooms
}

func (c *serverClient) IsInRoom(room string) bool {
	c.roomsMu.RLock()
	_, exists := c.rooms[room]
	c.roomsMu.RUnlock()
	return exists
}

// Authentication methods
func (c *serverClient) GetAuth() *AuthInfo {
	return c.auth
}

func (c *serverClient) SetAuth(auth *AuthInfo) {
	c.auth = auth
	if auth != nil {
		c.info.UserID = auth.UserID
	}
}

func (c *serverClient) IsAuthenticated() bool {
	return c.auth != nil
}

// Rate limiting methods
func (c *serverClient) IsRateLimited() bool {
	return !c.rateLimiter.Allow()
}

func (c *serverClient) GetRateLimit() (int, int) {
	// This is a simplified implementation
	// In a real implementation, you might want to track current usage
	return 0, c.options.RateLimit
}

// Connection management methods
func (c *serverClient) Ping() error {
	if !c.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	if c.options.WriteTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
	}

	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *serverClient) Close() error {
	return c.CloseWithCode(websocket.CloseNormalClosure, "")
}

func (c *serverClient) CloseWithCode(code int, text string) error {
	if !atomic.CompareAndSwapInt32(&c.state, int32(StateConnected), int32(StateDisconnected)) {
		return nil // already disconnected
	}

	// Send close message
	if c.options.WriteTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
	}
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))

	// Close connection
	c.conn.Close()
	c.cancel()
	close(c.closeCh)

	// Wait for goroutines to finish
	c.wg.Wait()

	return nil
}

// Context methods
func (c *serverClient) GetContext() context.Context {
	return c.ctx
}

func (c *serverClient) SetContext(ctx context.Context) {
	c.cancel()
	c.ctx, c.cancel = context.WithCancel(ctx)
}

// Custom data methods
func (c *serverClient) Set(key string, value any) {
	c.dataMu.Lock()
	c.data[key] = value
	c.dataMu.Unlock()
}

func (c *serverClient) Get(key string) (any, bool) {
	c.dataMu.RLock()
	value, exists := c.data[key]
	c.dataMu.RUnlock()
	return value, exists
}

func (c *serverClient) Delete(key string) {
	c.dataMu.Lock()
	delete(c.data, key)
	c.dataMu.Unlock()
}

// readPump đọc messages từ WebSocket connection
func (c *serverClient) readPump(server *wsServer) {
	defer func() {
		c.wg.Done()
		c.Close()
	}()

	// Set connection options
	if c.options.MaxMessageSize > 0 {
		c.conn.SetReadLimit(c.options.MaxMessageSize)
	}

	// Set read deadline
	if c.options.ReadTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.options.ReadTimeout))
	}

	// Set pong handler
	c.conn.SetPongHandler(func(string) error {
		c.pongMu.Lock()
		c.lastPong = time.Now()
		c.pongMu.Unlock()

		c.info.LastPing = time.Now()

		if c.options.ReadTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.options.ReadTimeout))
		}
		return nil
	})

	for {
		select {
		case <-c.closeCh:
			return
		case <-c.ctx.Done():
			return
		default:
		}

		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				server.handleError(c, fmt.Errorf("websocket read error: %w", err))
			}
			return
		}

		// Update metrics
		atomic.AddInt64(&c.metrics.MessagesReceived, 1)
		atomic.AddInt64(&c.metrics.BytesReceived, int64(len(data)))
		c.metrics.LastActivity = time.Now()

		// Handle message
		server.handleMessage(c, messageType, data)
	}
}

// writePump ghi messages đến WebSocket connection
func (c *serverClient) writePump(server *wsServer) {
	defer func() {
		c.wg.Done()
		c.conn.Close()
	}()

	ticker := time.NewTicker(c.options.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeCh:
			return
		case <-c.ctx.Done():
			return

		case msg := <-c.sendCh:
			// Set write deadline
			if c.options.WriteTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
			}

			var err error
			switch msg.Type {
			case TextMessage:
				err = c.conn.WriteMessage(websocket.TextMessage, msg.Data)
			case BinaryMessage:
				err = c.conn.WriteMessage(websocket.BinaryMessage, msg.Data)
			default:
				err = c.conn.WriteMessage(int(msg.Type), msg.Data)
			}

			if err != nil {
				server.handleError(c, fmt.Errorf("websocket write error: %w", err))
				return
			}

			// Update metrics
			atomic.AddInt64(&c.metrics.MessagesSent, 1)
			atomic.AddInt64(&c.metrics.BytesSent, int64(len(msg.Data)))
			c.metrics.LastActivity = time.Now()

		case <-ticker.C:
			// Send ping
			if c.options.WriteTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				server.handleError(c, fmt.Errorf("ping failed: %w", err))
				return
			}

			c.info.LastPing = time.Now()
		}
	}
}

// Helper functions
func getClientIP(req *http.Request) string {
	// Check X-Forwarded-For header
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Use remote address
	return req.RemoteAddr
}

func getHeaders(req *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	return headers
}
