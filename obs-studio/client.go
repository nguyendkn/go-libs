package obs_studio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents an OBS WebSocket client
type Client struct {
	config         *ConnectionConfig
	conn           *websocket.Conn
	authManager    *AuthManager
	requestManager *RequestManager
	eventProcessor *EventProcessor

	// Connection state
	connected  bool
	identified bool
	mutex      sync.RWMutex

	// Channels for internal communication
	sendChan  chan *WebSocketMessage
	closeChan chan struct{}

	// Context for connection management
	ctx    context.Context
	cancel context.CancelFunc

	// Connection info
	helloData      *HelloData
	identifiedData *IdentifiedData
}

// NewClient creates a new OBS WebSocket client
func NewClient() *Client {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		authManager:    NewAuthManager(),
		requestManager: NewRequestManager(30 * time.Second),
		eventProcessor: NewEventProcessor(),
		sendChan:       make(chan *WebSocketMessage, 100),
		closeChan:      make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}

	return client
}

// Connect establishes connection to OBS WebSocket server
func (c *Client) Connect(address, password string, opts ...ConnectOption) error {
	// Apply options
	config := &ConnectionConfig{
		Address:            address,
		Password:           password,
		RpcVersion:         1,
		EventSubscriptions: EventSubscriptionAll,
		ConnectTimeout:     10 * time.Second,
		RequestTimeout:     30 * time.Second,
	}

	for _, opt := range opts {
		opt(config)
	}

	c.config = config
	c.requestManager.timeout = config.RequestTimeout

	// Parse URL
	u, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Connect to WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: config.ConnectTimeout,
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.conn = conn
	c.connected = true

	// Start message handling
	go c.messageReader()
	go c.messageWriter()

	// Emit connection opened event
	c.eventProcessor.Emit(EventConnectionOpened, &ConnectionEventData{
		ConnectionStatus: "opened",
		Message:          "WebSocket connection established",
	})

	return nil
}

// Disconnect closes the WebSocket connection
func (c *Client) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	c.identified = false

	// Cancel context to stop goroutines
	c.cancel()

	// Close WebSocket connection
	if c.conn != nil {
		c.conn.Close()
	}

	// Cancel all pending requests
	c.requestManager.CancelAllRequests()

	// Emit connection closed event
	c.eventProcessor.Emit(EventConnectionClosed, &ConnectionEventData{
		ConnectionStatus: "closed",
		Message:          "WebSocket connection closed",
	})

	return nil
}

// IsConnected returns true if client is connected
func (c *Client) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.connected
}

// IsIdentified returns true if client is identified
func (c *Client) IsIdentified() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.identified
}

// Call sends a request and waits for response
func (c *Client) Call(requestType string, requestData interface{}) (*RequestResponseData, error) {
	if !c.IsIdentified() {
		return nil, fmt.Errorf("client must be identified before making requests")
	}

	request := c.requestManager.CreateRequest(requestType, requestData)
	return c.requestManager.SendRequest(c.ctx, request, c.sendMessage)
}

// CallBatch sends a batch of requests
func (c *Client) CallBatch(requests []RequestData, haltOnFailure bool) (*BatchResponseData, error) {
	if !c.IsIdentified() {
		return nil, fmt.Errorf("client must be identified before making requests")
	}

	return c.requestManager.SendBatchRequest(c.ctx, requests, haltOnFailure, c.sendMessage)
}

// On registers an event handler
func (c *Client) On(eventType string, handler EventHandler) {
	c.eventProcessor.On(eventType, handler)
}

// Off removes an event handler
func (c *Client) Off(eventType string, handler EventHandler) {
	c.eventProcessor.Off(eventType, handler)
}

// messageReader handles incoming WebSocket messages
func (c *Client) messageReader() {
	defer func() {
		if r := recover(); r != nil {
			c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
				ConnectionStatus: "error",
				Message:          fmt.Sprintf("Message reader panic: %v", r),
			})
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			messageType, data, err := c.conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
						ConnectionStatus: "error",
						Message:          fmt.Sprintf("Read error: %v", err),
					})
				}
				return
			}

			if messageType != websocket.TextMessage {
				continue
			}

			message := &WebSocketMessage{}
			if err := json.Unmarshal(data, message); err != nil {
				c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
					ConnectionStatus: "error",
					Message:          fmt.Sprintf("JSON unmarshal error: %v", err),
				})
				continue
			}

			c.processMessage(message)
		}
	}
}

// messageWriter handles outgoing WebSocket messages
func (c *Client) messageWriter() {
	defer func() {
		if r := recover(); r != nil {
			c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
				ConnectionStatus: "error",
				Message:          fmt.Sprintf("Message writer panic: %v", r),
			})
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.sendChan:
			data, err := json.Marshal(message)
			if err != nil {
				c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
					ConnectionStatus: "error",
					Message:          fmt.Sprintf("JSON marshal error: %v", err),
				})
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
					ConnectionStatus: "error",
					Message:          fmt.Sprintf("Write error: %v", err),
				})
				return
			}
		}
	}
}

// processMessage processes incoming WebSocket messages
func (c *Client) processMessage(message *WebSocketMessage) {
	switch message.Op {
	case OpCodeHello:
		c.handleHello(message)
	case OpCodeIdentified:
		c.handleIdentified(message)
	case OpCodeEvent:
		c.eventProcessor.ProcessEvent(message)
	case OpCodeRequestResponse, OpCodeRequestBatchResponse:
		c.requestManager.ProcessResponse(message)
	default:
		c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
			ConnectionStatus: "error",
			Message:          fmt.Sprintf("Unknown op code: %d", message.Op),
		})
	}
}

// handleHello handles Hello message from server
func (c *Client) handleHello(message *WebSocketMessage) {
	helloData := &HelloData{}
	dataBytes, err := json.Marshal(message.D)
	if err != nil {
		c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
			ConnectionStatus: "error",
			Message:          fmt.Sprintf("Failed to parse hello data: %v", err),
		})
		return
	}

	if err := json.Unmarshal(dataBytes, helloData); err != nil {
		c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
			ConnectionStatus: "error",
			Message:          fmt.Sprintf("Failed to unmarshal hello data: %v", err),
		})
		return
	}

	c.helloData = helloData
	c.sendIdentify()
}

// sendIdentify sends Identify message to server
func (c *Client) sendIdentify() {
	identifyData := &IdentifyData{
		RpcVersion:         c.config.RpcVersion,
		EventSubscriptions: c.config.EventSubscriptions,
	}

	// Handle authentication if required
	if c.authManager.RequiresAuthentication(c.helloData) {
		if c.config.Password == "" {
			c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
				ConnectionStatus: "error",
				Message:          "Authentication required but no password provided",
			})
			return
		}

		authString, err := c.authManager.GenerateAuthString(c.config.Password, c.helloData.Authentication)
		if err != nil {
			c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
				ConnectionStatus: "error",
				Message:          fmt.Sprintf("Authentication failed: %v", err),
			})
			return
		}

		identifyData.Authentication = authString
	}

	message := &WebSocketMessage{
		Op: OpCodeIdentify,
		D:  identifyData,
	}

	c.sendMessage(message)
}

// handleIdentified handles Identified message from server
func (c *Client) handleIdentified(message *WebSocketMessage) {
	identifiedData := &IdentifiedData{}
	dataBytes, err := json.Marshal(message.D)
	if err != nil {
		c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
			ConnectionStatus: "error",
			Message:          fmt.Sprintf("Failed to parse identified data: %v", err),
		})
		return
	}

	if err := json.Unmarshal(dataBytes, identifiedData); err != nil {
		c.eventProcessor.Emit(EventConnectionError, &ConnectionEventData{
			ConnectionStatus: "error",
			Message:          fmt.Sprintf("Failed to unmarshal identified data: %v", err),
		})
		return
	}

	c.mutex.Lock()
	c.identified = true
	c.identifiedData = identifiedData
	c.mutex.Unlock()

	// Emit identified event
	c.eventProcessor.Emit(EventIdentified, identifiedData)
}

// sendMessage sends a message through the send channel
func (c *Client) sendMessage(message *WebSocketMessage) error {
	select {
	case c.sendChan <- message:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("client is closed")
	default:
		return fmt.Errorf("send channel is full")
	}
}

// ConnectOption represents connection configuration option
type ConnectOption func(*ConnectionConfig)

// WithRpcVersion sets the RPC version
func WithRpcVersion(version int) ConnectOption {
	return func(config *ConnectionConfig) {
		config.RpcVersion = version
	}
}

// WithEventSubscriptions sets event subscription flags
func WithEventSubscriptions(subscriptions int) ConnectOption {
	return func(config *ConnectionConfig) {
		config.EventSubscriptions = subscriptions
	}
}

// WithConnectTimeout sets connection timeout
func WithConnectTimeout(timeout time.Duration) ConnectOption {
	return func(config *ConnectionConfig) {
		config.ConnectTimeout = timeout
	}
}

// WithRequestTimeout sets request timeout
func WithRequestTimeout(timeout time.Duration) ConnectOption {
	return func(config *ConnectionConfig) {
		config.RequestTimeout = timeout
	}
}

// Helper methods for common operations

// GetVersion gets OBS version information
func (c *Client) GetVersion() (*RequestResponseData, error) {
	return c.Call(RequestTypeGetVersion, nil)
}

// GetSceneList gets list of scenes
func (c *Client) GetSceneList() (*SceneListResponse, error) {
	response, err := c.Call(RequestTypeGetSceneList, nil)
	if err != nil {
		return nil, err
	}

	sceneList := &SceneListResponse{}
	responseBytes, err := json.Marshal(response.ResponseData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(responseBytes, sceneList); err != nil {
		return nil, err
	}

	return sceneList, nil
}

// SetCurrentScene sets the current program scene
func (c *Client) SetCurrentScene(sceneName string) error {
	_, err := c.Call(RequestTypeSetCurrentProgramScene, &SetCurrentSceneRequestData{
		SceneName: sceneName,
	})
	return err
}

// StartStream starts streaming
func (c *Client) StartStream() error {
	_, err := c.Call(RequestTypeStartStream, nil)
	return err
}

// StopStream stops streaming
func (c *Client) StopStream() error {
	_, err := c.Call(RequestTypeStopStream, nil)
	return err
}

// StartRecord starts recording
func (c *Client) StartRecord() error {
	_, err := c.Call(RequestTypeStartRecord, nil)
	return err
}

// StopRecord stops recording
func (c *Client) StopRecord() error {
	_, err := c.Call(RequestTypeStopRecord, nil)
	return err
}

// GetRequestManager returns the internal request manager for advanced usage
func (c *Client) GetRequestManager() *RequestManager {
	return c.requestManager
}
