package webrtc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// signalingClient implements the SignalingClient interface
type signalingClient struct {
	// Connection
	conn      *websocket.Conn
	url       string
	connected int32 // atomic

	// Authentication
	authToken string

	// Event handlers
	onMessage      func(*SignalingMessage)
	onOffer        func(string, *SessionDescription)
	onAnswer       func(string, *SessionDescription)
	onICECandidate func(string, *ICECandidate)
	onPeerJoined   func(*PeerInfo)
	onPeerLeft     func(string)
	onRoomUpdate   func(*RoomInfo)
	onError        func(error)
	handlersMu     sync.RWMutex

	// Reconnection
	reconnectEnabled     bool
	reconnectInterval    time.Duration
	maxReconnectAttempts int
	reconnectAttempts    int

	// Context and lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Channels
	sendCh chan *SignalingMessage

	// Wait group
	wg sync.WaitGroup
}

// NewSignalingClient tạo một SignalingClient mới
func NewSignalingClient() SignalingClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &signalingClient{
		sendCh:               make(chan *SignalingMessage, 100),
		reconnectEnabled:     true,
		reconnectInterval:    5 * time.Second,
		maxReconnectAttempts: 10,
		ctx:                  ctx,
		cancel:               cancel,
	}
}

// Connection management
func (sc *signalingClient) Connect(url string) error {
	return sc.ConnectWithContext(context.Background(), url)
}

func (sc *signalingClient) ConnectWithContext(ctx context.Context, urlStr string) error {
	if atomic.LoadInt32(&sc.connected) == 1 {
		return fmt.Errorf("already connected")
	}

	sc.url = urlStr

	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Add auth token to headers if available
	headers := make(map[string][]string)
	if sc.authToken != "" {
		headers["Authorization"] = []string{"Bearer " + sc.authToken}
	}

	// Connect to WebSocket
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	sc.conn = conn
	atomic.StoreInt32(&sc.connected, 1)
	sc.reconnectAttempts = 0

	// Start goroutines
	sc.wg.Add(2)
	go sc.readPump()
	go sc.writePump()

	return nil
}

func (sc *signalingClient) Disconnect() error {
	if !atomic.CompareAndSwapInt32(&sc.connected, 1, 0) {
		return nil // already disconnected
	}

	sc.cancel()

	if sc.conn != nil {
		sc.conn.Close()
	}

	sc.wg.Wait()

	return nil
}

func (sc *signalingClient) IsConnected() bool {
	return atomic.LoadInt32(&sc.connected) == 1
}

// Message handling
func (sc *signalingClient) SendMessage(msg *SignalingMessage) error {
	if !sc.IsConnected() {
		return fmt.Errorf("not connected")
	}

	msg.Timestamp = time.Now()

	select {
	case sc.sendCh <- msg:
		return nil
	case <-sc.ctx.Done():
		return fmt.Errorf("context cancelled")
	default:
		return fmt.Errorf("send queue is full")
	}
}

func (sc *signalingClient) SendOffer(to string, offer *SessionDescription) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "offer",
		To:   to,
		Data: offer,
	})
}

func (sc *signalingClient) SendAnswer(to string, answer *SessionDescription) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "answer",
		To:   to,
		Data: answer,
	})
}

func (sc *signalingClient) SendICECandidate(to string, candidate *ICECandidate) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "ice-candidate",
		To:   to,
		Data: candidate,
	})
}

func (sc *signalingClient) SendBye(to string) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "bye",
		To:   to,
	})
}

// Room management
func (sc *signalingClient) JoinRoom(roomID string, userInfo *PeerInfo) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "join-room",
		Room: roomID,
		Data: userInfo,
	})
}

func (sc *signalingClient) LeaveRoom(roomID string) error {
	return sc.SendMessage(&SignalingMessage{
		Type: "leave-room",
		Room: roomID,
	})
}

func (sc *signalingClient) GetRoomInfo(roomID string) (*RoomInfo, error) {
	// This would typically be implemented as a request-response pattern
	// For now, return an error indicating it's not implemented
	return nil, fmt.Errorf("GetRoomInfo not implemented - use async events instead")
}

func (sc *signalingClient) ListRooms() ([]*RoomInfo, error) {
	// This would typically be implemented as a request-response pattern
	// For now, return an error indicating it's not implemented
	return nil, fmt.Errorf("ListRooms not implemented - use async events instead")
}

// Event handlers
func (sc *signalingClient) OnMessage(handler func(*SignalingMessage)) {
	sc.handlersMu.Lock()
	sc.onMessage = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnOffer(handler func(string, *SessionDescription)) {
	sc.handlersMu.Lock()
	sc.onOffer = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnAnswer(handler func(string, *SessionDescription)) {
	sc.handlersMu.Lock()
	sc.onAnswer = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnICECandidate(handler func(string, *ICECandidate)) {
	sc.handlersMu.Lock()
	sc.onICECandidate = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnPeerJoined(handler func(*PeerInfo)) {
	sc.handlersMu.Lock()
	sc.onPeerJoined = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnPeerLeft(handler func(string)) {
	sc.handlersMu.Lock()
	sc.onPeerLeft = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnRoomUpdate(handler func(*RoomInfo)) {
	sc.handlersMu.Lock()
	sc.onRoomUpdate = handler
	sc.handlersMu.Unlock()
}

func (sc *signalingClient) OnError(handler func(error)) {
	sc.handlersMu.Lock()
	sc.onError = handler
	sc.handlersMu.Unlock()
}

// Authentication
func (sc *signalingClient) SetAuthToken(token string) {
	sc.authToken = token
}

func (sc *signalingClient) GetAuthToken() string {
	return sc.authToken
}

// Configuration
func (sc *signalingClient) SetReconnectOptions(enabled bool, interval, maxAttempts int) {
	sc.reconnectEnabled = enabled
	sc.reconnectInterval = time.Duration(interval) * time.Second
	sc.maxReconnectAttempts = maxAttempts
}

// Lifecycle
func (sc *signalingClient) Close() error {
	return sc.Disconnect()
}

// readPump đọc messages từ WebSocket
func (sc *signalingClient) readPump() {
	defer func() {
		sc.wg.Done()
		sc.handleDisconnection()
	}()

	for {
		select {
		case <-sc.ctx.Done():
			return
		default:
		}

		_, data, err := sc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sc.emitError(fmt.Errorf("websocket read error: %w", err))
			}
			return
		}

		var msg SignalingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			sc.emitError(fmt.Errorf("failed to unmarshal message: %w", err))
			continue
		}

		sc.handleMessage(&msg)
	}
}

// writePump ghi messages đến WebSocket
func (sc *signalingClient) writePump() {
	defer func() {
		sc.wg.Done()
		sc.conn.Close()
	}()

	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sc.ctx.Done():
			return

		case msg := <-sc.sendCh:
			sc.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			data, err := json.Marshal(msg)
			if err != nil {
				sc.emitError(fmt.Errorf("failed to marshal message: %w", err))
				continue
			}

			if err := sc.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				sc.emitError(fmt.Errorf("websocket write error: %w", err))
				return
			}

		case <-ticker.C:
			sc.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := sc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage xử lý message nhận được
func (sc *signalingClient) handleMessage(msg *SignalingMessage) {
	sc.handlersMu.RLock()
	defer sc.handlersMu.RUnlock()

	// Call general message handler
	if sc.onMessage != nil {
		go sc.onMessage(msg)
	}

	// Call specific handlers based on message type
	switch msg.Type {
	case "offer":
		if sc.onOffer != nil && msg.Data != nil {
			if offer, ok := msg.Data.(*SessionDescription); ok {
				go sc.onOffer(msg.From, offer)
			}
		}

	case "answer":
		if sc.onAnswer != nil && msg.Data != nil {
			if answer, ok := msg.Data.(*SessionDescription); ok {
				go sc.onAnswer(msg.From, answer)
			}
		}

	case "ice-candidate":
		if sc.onICECandidate != nil && msg.Data != nil {
			if candidate, ok := msg.Data.(*ICECandidate); ok {
				go sc.onICECandidate(msg.From, candidate)
			}
		}

	case "peer-joined":
		if sc.onPeerJoined != nil && msg.Data != nil {
			if peerInfo, ok := msg.Data.(*PeerInfo); ok {
				go sc.onPeerJoined(peerInfo)
			}
		}

	case "peer-left":
		if sc.onPeerLeft != nil {
			go sc.onPeerLeft(msg.From)
		}

	case "room-update":
		if sc.onRoomUpdate != nil && msg.Data != nil {
			if roomInfo, ok := msg.Data.(*RoomInfo); ok {
				go sc.onRoomUpdate(roomInfo)
			}
		}
	}
}

// handleDisconnection xử lý khi mất kết nối
func (sc *signalingClient) handleDisconnection() {
	atomic.StoreInt32(&sc.connected, 0)

	if sc.reconnectEnabled && sc.reconnectAttempts < sc.maxReconnectAttempts {
		sc.reconnectAttempts++

		time.Sleep(sc.reconnectInterval)

		if err := sc.Connect(sc.url); err != nil {
			sc.emitError(fmt.Errorf("reconnection attempt %d failed: %w", sc.reconnectAttempts, err))
		}
	}
}

// emitError emit error event
func (sc *signalingClient) emitError(err error) {
	sc.handlersMu.RLock()
	if sc.onError != nil {
		go sc.onError(err)
	}
	sc.handlersMu.RUnlock()
}
