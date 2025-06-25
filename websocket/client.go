package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nguyendkn/go-libs/uuid"
	"golang.org/x/time/rate"
)

// wsClient implements the Client interface
type wsClient struct {
	id      string
	options *ClientOptions

	// Connection
	conn   *websocket.Conn
	connMu sync.RWMutex
	state  int32 // atomic, ConnectionState

	// Context and cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// Channels
	sendCh      chan *Message
	closeCh     chan struct{}
	reconnectCh chan struct{}

	// Event handlers
	onConnect       func()
	onDisconnect    func(error)
	onMessage       func([]byte)
	onTextMessage   func(string)
	onBinaryMessage func([]byte)
	onError         func(error)
	onReconnect     func(int)
	handlersMu      sync.RWMutex

	// Metrics and info
	info    *ClientInfo
	metrics *ConnectionMetrics

	// Rate limiter
	rateLimiter *rate.Limiter

	// Reconnection
	reconnectAttempts int32
	lastReconnect     time.Time

	// Heartbeat
	lastPong time.Time
	pongMu   sync.RWMutex

	// Wait groups
	wg sync.WaitGroup
}

// NewClient tạo một WebSocket client mới
func NewClient(url string, options ...*ClientOptions) Client {
	opts := &ClientOptions{
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
	}

	if len(options) > 0 && options[0] != nil {
		mergeClientOptions(opts, options[0])
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &wsClient{
		id:          uuid.New().String(),
		options:     opts,
		ctx:         ctx,
		cancel:      cancel,
		sendCh:      make(chan *Message, opts.MessageQueueSize),
		closeCh:     make(chan struct{}),
		reconnectCh: make(chan struct{}, 1),
		info: &ClientInfo{
			ID:          uuid.New().String(),
			ConnectedAt: time.Now(),
		},
		metrics:     &ConnectionMetrics{},
		rateLimiter: rate.NewLimiter(rate.Limit(opts.RateLimit), opts.RateBurst),
	}

	atomic.StoreInt32(&client.state, int32(StateDisconnected))

	return client
}

// Connect kết nối đến WebSocket server
func (c *wsClient) Connect() error {
	return c.ConnectWithContext(context.Background())
}

// ConnectWithContext kết nối với context
func (c *wsClient) ConnectWithContext(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&c.state, int32(StateDisconnected), int32(StateConnecting)) {
		return fmt.Errorf("client is already connecting or connected")
	}

	// Parse URL
	u, err := url.Parse(c.options.URL)
	if err != nil {
		atomic.StoreInt32(&c.state, int32(StateDisconnected))
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare headers
	headers := http.Header{}
	for k, v := range c.options.Headers {
		headers.Set(k, v)
	}

	// Add auth header if provided
	if c.options.AuthToken != "" {
		authHeader := c.options.AuthHeader
		if authHeader == "" {
			authHeader = "Authorization"
		}
		headers.Set(authHeader, c.options.AuthToken)
	}

	// Get auth token from callback if provided
	if c.options.AuthCallback != nil {
		token, err := c.options.AuthCallback()
		if err != nil {
			atomic.StoreInt32(&c.state, int32(StateDisconnected))
			return fmt.Errorf("auth callback failed: %w", err)
		}
		if token != "" {
			authHeader := c.options.AuthHeader
			if authHeader == "" {
				authHeader = "Authorization"
			}
			headers.Set(authHeader, token)
		}
	}

	// Create dialer
	dialer := c.options.Dialer
	if dialer == nil {
		dialer = &websocket.Dialer{
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: c.options.EnableCompression,
		}
	}

	// Set proxy if provided
	if c.options.ProxyURL != "" {
		proxyURL, err := url.Parse(c.options.ProxyURL)
		if err != nil {
			atomic.StoreInt32(&c.state, int32(StateDisconnected))
			return fmt.Errorf("invalid proxy URL: %w", err)
		}
		dialer.Proxy = http.ProxyURL(proxyURL)
	}

	// Connect
	conn, resp, err := dialer.DialContext(ctx, u.String(), headers)
	if err != nil {
		atomic.StoreInt32(&c.state, int32(StateDisconnected))
		if resp != nil {
			return fmt.Errorf("websocket dial failed (status %d): %w", resp.StatusCode, err)
		}
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	// Set connection options
	if c.options.MaxMessageSize > 0 {
		conn.SetReadLimit(c.options.MaxMessageSize)
	}

	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	// Update info
	c.info.ConnectedAt = time.Now()
	c.info.LastPing = time.Now()
	c.pongMu.Lock()
	c.lastPong = time.Now()
	c.pongMu.Unlock()

	atomic.StoreInt32(&c.state, int32(StateConnected))
	atomic.StoreInt32(&c.reconnectAttempts, 0)

	// Start goroutines
	c.wg.Add(3)
	go c.readPump()
	go c.writePump()
	go c.heartbeat()

	// Call connect handler
	c.handlersMu.RLock()
	if c.onConnect != nil {
		go c.onConnect()
	}
	c.handlersMu.RUnlock()

	return nil
}

// Disconnect ngắt kết nối
func (c *wsClient) Disconnect() error {
	if !atomic.CompareAndSwapInt32(&c.state, int32(StateConnected), int32(StateDisconnected)) &&
		!atomic.CompareAndSwapInt32(&c.state, int32(StateReconnecting), int32(StateDisconnected)) {
		return nil // already disconnected
	}

	c.connMu.Lock()
	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
		c.conn = nil
	}
	c.connMu.Unlock()

	close(c.closeCh)
	c.wg.Wait()

	return nil
}

// IsConnected kiểm tra trạng thái kết nối
func (c *wsClient) IsConnected() bool {
	return atomic.LoadInt32(&c.state) == int32(StateConnected)
}

// GetState trả về trạng thái hiện tại
func (c *wsClient) GetState() ConnectionState {
	return ConnectionState(atomic.LoadInt32(&c.state))
}

// Send gửi message binary
func (c *wsClient) Send(data []byte) error {
	return c.SendWithType(BinaryMessage, data)
}

// SendText gửi message text
func (c *wsClient) SendText(text string) error {
	return c.SendWithType(TextMessage, []byte(text))
}

// SendJSON gửi message JSON
func (c *wsClient) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return c.SendWithType(TextMessage, data)
}

// SendWithType gửi message với type cụ thể
func (c *wsClient) SendWithType(messageType MessageType, data []byte) error {
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

// Event handler setters
func (c *wsClient) OnConnect(handler func()) {
	c.handlersMu.Lock()
	c.onConnect = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnDisconnect(handler func(error)) {
	c.handlersMu.Lock()
	c.onDisconnect = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnMessage(handler func([]byte)) {
	c.handlersMu.Lock()
	c.onMessage = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnTextMessage(handler func(string)) {
	c.handlersMu.Lock()
	c.onTextMessage = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnBinaryMessage(handler func([]byte)) {
	c.handlersMu.Lock()
	c.onBinaryMessage = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnError(handler func(error)) {
	c.handlersMu.Lock()
	c.onError = handler
	c.handlersMu.Unlock()
}

func (c *wsClient) OnReconnect(handler func(int)) {
	c.handlersMu.Lock()
	c.onReconnect = handler
	c.handlersMu.Unlock()
}

// ID trả về client ID
func (c *wsClient) ID() string {
	return c.id
}

// Info trả về thông tin client
func (c *wsClient) Info() *ClientInfo {
	return c.info
}

// Metrics trả về metrics
func (c *wsClient) Metrics() *ConnectionMetrics {
	return c.metrics
}

// SetOptions cập nhật options
func (c *wsClient) SetOptions(options *ClientOptions) {
	if options != nil {
		mergeClientOptions(c.options, options)
		// Update rate limiter if needed
		if options.RateLimit > 0 {
			c.rateLimiter = rate.NewLimiter(rate.Limit(options.RateLimit), options.RateBurst)
		}
	}
}

// GetOptions trả về options hiện tại
func (c *wsClient) GetOptions() *ClientOptions {
	return c.options
}

// Close đóng client
func (c *wsClient) Close() error {
	atomic.StoreInt32(&c.state, int32(StateClosed))
	c.cancel()
	return c.Disconnect()
}

// readPump đọc messages từ WebSocket connection
func (c *wsClient) readPump() {
	defer c.wg.Done()

	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return
	}

	// Set read deadline
	if c.options.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(c.options.ReadTimeout))
	}

	// Set pong handler
	conn.SetPongHandler(func(string) error {
		c.pongMu.Lock()
		c.lastPong = time.Now()
		c.pongMu.Unlock()

		if c.options.ReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(c.options.ReadTimeout))
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

		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.handleError(fmt.Errorf("websocket read error: %w", err))
			}
			c.handleDisconnect(err)
			return
		}

		// Update metrics
		atomic.AddInt64(&c.metrics.MessagesReceived, 1)
		atomic.AddInt64(&c.metrics.BytesReceived, int64(len(data)))
		c.metrics.LastActivity = time.Now()

		// Handle message based on type
		switch messageType {
		case websocket.TextMessage:
			c.handlersMu.RLock()
			if c.onTextMessage != nil {
				go c.onTextMessage(string(data))
			}
			if c.onMessage != nil {
				go c.onMessage(data)
			}
			c.handlersMu.RUnlock()

		case websocket.BinaryMessage:
			c.handlersMu.RLock()
			if c.onBinaryMessage != nil {
				go c.onBinaryMessage(data)
			}
			if c.onMessage != nil {
				go c.onMessage(data)
			}
			c.handlersMu.RUnlock()

		case websocket.CloseMessage:
			c.handleDisconnect(nil)
			return
		}
	}
}

// writePump ghi messages đến WebSocket connection
func (c *wsClient) writePump() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.options.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeCh:
			return
		case <-c.ctx.Done():
			return

		case msg := <-c.sendCh:
			c.connMu.RLock()
			conn := c.conn
			c.connMu.RUnlock()

			if conn == nil {
				continue
			}

			// Set write deadline
			if c.options.WriteTimeout > 0 {
				conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
			}

			var err error
			switch msg.Type {
			case TextMessage:
				err = conn.WriteMessage(websocket.TextMessage, msg.Data)
			case BinaryMessage:
				err = conn.WriteMessage(websocket.BinaryMessage, msg.Data)
			default:
				err = conn.WriteMessage(int(msg.Type), msg.Data)
			}

			if err != nil {
				c.handleError(fmt.Errorf("websocket write error: %w", err))
				c.handleDisconnect(err)
				return
			}

			// Update metrics
			atomic.AddInt64(&c.metrics.MessagesSent, 1)
			atomic.AddInt64(&c.metrics.BytesSent, int64(len(msg.Data)))
			c.metrics.LastActivity = time.Now()

		case <-ticker.C:
			c.connMu.RLock()
			conn := c.conn
			c.connMu.RUnlock()

			if conn == nil {
				continue
			}

			// Send ping
			if c.options.WriteTimeout > 0 {
				conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
			}

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.handleError(fmt.Errorf("ping failed: %w", err))
				c.handleDisconnect(err)
				return
			}

			c.info.LastPing = time.Now()
		}
	}
}

// heartbeat kiểm tra heartbeat
func (c *wsClient) heartbeat() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.options.PongTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeCh:
			return
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.pongMu.RLock()
			lastPong := c.lastPong
			c.pongMu.RUnlock()

			if time.Since(lastPong) > c.options.PongTimeout {
				c.handleError(fmt.Errorf("pong timeout"))
				c.handleDisconnect(fmt.Errorf("heartbeat timeout"))
				return
			}
		}
	}
}

// handleDisconnect xử lý disconnect
func (c *wsClient) handleDisconnect(err error) {
	if !atomic.CompareAndSwapInt32(&c.state, int32(StateConnected), int32(StateDisconnected)) {
		return // already disconnected
	}

	c.connMu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connMu.Unlock()

	// Call disconnect handler
	c.handlersMu.RLock()
	if c.onDisconnect != nil {
		go c.onDisconnect(err)
	}
	c.handlersMu.RUnlock()

	// Try to reconnect if enabled
	if c.options.AutoReconnect && atomic.LoadInt32(&c.state) != int32(StateClosed) {
		go c.reconnect()
	}
}

// handleError xử lý error
func (c *wsClient) handleError(err error) {
	atomic.AddInt64(&c.metrics.Errors, 1)

	c.handlersMu.RLock()
	if c.onError != nil {
		go c.onError(err)
	}
	c.handlersMu.RUnlock()
}

// reconnect thực hiện reconnect
func (c *wsClient) reconnect() {
	attempts := atomic.LoadInt32(&c.reconnectAttempts)
	if c.options.MaxReconnectAttempts > 0 && int(attempts) >= c.options.MaxReconnectAttempts {
		atomic.StoreInt32(&c.state, int32(StateDisconnected))
		return
	}

	atomic.StoreInt32(&c.state, int32(StateReconnecting))
	atomic.AddInt32(&c.reconnectAttempts, 1)

	// Calculate backoff delay
	delay := c.options.ReconnectInterval
	if c.options.ReconnectBackoff > 0 {
		backoffMultiplier := time.Duration(attempts)
		if backoffMultiplier > 10 {
			backoffMultiplier = 10 // cap at 10x
		}
		delay = c.options.ReconnectInterval + (c.options.ReconnectBackoff * backoffMultiplier)
	}

	// Wait before reconnecting
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-c.ctx.Done():
		return
	}

	// Call reconnect handler
	c.handlersMu.RLock()
	if c.onReconnect != nil {
		go c.onReconnect(int(attempts))
	}
	c.handlersMu.RUnlock()

	// Try to connect
	if err := c.ConnectWithContext(c.ctx); err != nil {
		// Reconnect failed, try again
		go c.reconnect()
	}
}

// mergeClientOptions merge options
func mergeClientOptions(dst, src *ClientOptions) {
	if src.URL != "" {
		dst.URL = src.URL
	}
	if src.Headers != nil {
		if dst.Headers == nil {
			dst.Headers = make(map[string]string)
		}
		for k, v := range src.Headers {
			dst.Headers[k] = v
		}
	}
	if src.Subprotocols != nil {
		dst.Subprotocols = src.Subprotocols
	}
	if src.ReconnectInterval > 0 {
		dst.ReconnectInterval = src.ReconnectInterval
	}
	if src.MaxReconnectAttempts > 0 {
		dst.MaxReconnectAttempts = src.MaxReconnectAttempts
	}
	if src.ReconnectBackoff > 0 {
		dst.ReconnectBackoff = src.ReconnectBackoff
	}
	if src.PingInterval > 0 {
		dst.PingInterval = src.PingInterval
	}
	if src.PongTimeout > 0 {
		dst.PongTimeout = src.PongTimeout
	}
	if src.WriteTimeout > 0 {
		dst.WriteTimeout = src.WriteTimeout
	}
	if src.ReadTimeout > 0 {
		dst.ReadTimeout = src.ReadTimeout
	}
	if src.MessageQueueSize > 0 {
		dst.MessageQueueSize = src.MessageQueueSize
	}
	if src.MaxMessageSize > 0 {
		dst.MaxMessageSize = src.MaxMessageSize
	}
	if src.RateLimit > 0 {
		dst.RateLimit = src.RateLimit
	}
	if src.RateBurst > 0 {
		dst.RateBurst = src.RateBurst
	}
	if src.RateWindow > 0 {
		dst.RateWindow = src.RateWindow
	}
	if src.AuthToken != "" {
		dst.AuthToken = src.AuthToken
	}
	if src.AuthHeader != "" {
		dst.AuthHeader = src.AuthHeader
	}
	if src.AuthCallback != nil {
		dst.AuthCallback = src.AuthCallback
	}
	if src.ProxyURL != "" {
		dst.ProxyURL = src.ProxyURL
	}
	if src.Dialer != nil {
		dst.Dialer = src.Dialer
	}
	dst.AutoReconnect = src.AutoReconnect
	dst.EnableCompression = src.EnableCompression
}
