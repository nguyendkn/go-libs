package webrtc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// signalingServer implements the SignalingServer interface
type signalingServer struct {
	// HTTP server
	server *http.Server

	// WebSocket upgrader
	upgrader websocket.Upgrader

	// Configuration
	config *ServerConfig

	// Room management
	rooms   map[string]Room
	roomsMu sync.RWMutex

	// Peer management
	peers   map[string]*signalingPeer
	peersMu sync.RWMutex

	// Event handlers
	onPeerConnected    func(*PeerInfo)
	onPeerDisconnected func(string)
	onRoomCreated      func(*RoomInfo)
	onRoomDeleted      func(string)
	onMessage          func(*SignalingMessage)
	onError            func(error)
	handlersMu         sync.RWMutex

	// Authentication
	authHandler func(string) (*PeerInfo, error)

	// Middleware
	middlewares []SignalingMiddleware

	// Statistics
	stats     *ServerStats
	statsMu   sync.RWMutex
	startTime time.Time

	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

// signalingPeer đại diện cho một peer kết nối
type signalingPeer struct {
	info *PeerInfo
	conn *websocket.Conn
	send chan *SignalingMessage

	server *signalingServer

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSignalingServer tạo một SignalingServer mới
func NewSignalingServer() SignalingServer {
	ctx, cancel := context.WithCancel(context.Background())

	server := &signalingServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins by default
			},
		},
		config: &ServerConfig{
			MaxRooms:        1000,
			MaxPeersPerRoom: 50,
			MessageTimeout:  30 * time.Second,
			PeerTimeout:     60 * time.Second,
			EnableAuth:      false,
			EnableCORS:      true,
			AllowedOrigins:  []string{"*"},
		},
		rooms:     make(map[string]Room),
		peers:     make(map[string]*signalingPeer),
		stats:     &ServerStats{},
		startTime: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	return server
}

// Server lifecycle
func (ss *signalingServer) Start(addr string) error {
	return ss.StartWithContext(context.Background(), addr)
}

func (ss *signalingServer) StartWithContext(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ss.handleWebSocket)
	mux.HandleFunc("/rooms", ss.handleRooms)
	mux.HandleFunc("/stats", ss.handleStats)

	ss.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go ss.statsCollector()

	return ss.server.ListenAndServe()
}

func (ss *signalingServer) Stop() error {
	return ss.Shutdown(context.Background())
}

func (ss *signalingServer) Shutdown(ctx context.Context) error {
	ss.cancel()

	// Close all peer connections
	ss.peersMu.Lock()
	for _, peer := range ss.peers {
		peer.cancel()
		peer.conn.Close()
	}
	ss.peersMu.Unlock()

	// Close all rooms
	ss.roomsMu.Lock()
	for _, room := range ss.rooms {
		room.Close()
	}
	ss.roomsMu.Unlock()

	if ss.server != nil {
		return ss.server.Shutdown(ctx)
	}

	return nil
}

// handleWebSocket xử lý WebSocket connections
func (ss *signalingServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := ss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ss.emitError(fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}

	// Authenticate if required
	var peerInfo *PeerInfo
	if ss.config.EnableAuth && ss.authHandler != nil {
		token := r.Header.Get("Authorization")
		if token == "" {
			conn.Close()
			return
		}

		// Remove "Bearer " prefix
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		peerInfo, err = ss.authHandler(token)
		if err != nil {
			conn.Close()
			return
		}
	} else {
		// Create anonymous peer info
		peerInfo = &PeerInfo{
			ID:       generatePeerID(),
			Username: "anonymous",
			JoinedAt: time.Now(),
		}
	}

	// Create peer
	ctx, cancel := context.WithCancel(ss.ctx)
	peer := &signalingPeer{
		info:   peerInfo,
		conn:   conn,
		send:   make(chan *SignalingMessage, 100),
		server: ss,
		ctx:    ctx,
		cancel: cancel,
	}

	// Register peer
	ss.peersMu.Lock()
	ss.peers[peerInfo.ID] = peer
	ss.statsMu.Lock()
	ss.stats.ActiveConnections++
	ss.statsMu.Unlock()
	ss.peersMu.Unlock()

	// Emit peer connected event
	ss.handlersMu.RLock()
	if ss.onPeerConnected != nil {
		go ss.onPeerConnected(peerInfo)
	}
	ss.handlersMu.RUnlock()

	// Start peer goroutines
	go peer.readPump()
	go peer.writePump()
}

// handleRooms xử lý REST API cho rooms
func (ss *signalingServer) handleRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		rooms, err := ss.ListRooms()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(rooms)

	case "POST":
		var roomInfo RoomInfo
		if err := json.NewDecoder(r.Body).Decode(&roomInfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ss.CreateRoom(&roomInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(roomInfo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleStats xử lý stats endpoint
func (ss *signalingServer) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ss.GetStats())
}

// Room management
func (ss *signalingServer) CreateRoom(info *RoomInfo) error {
	ss.roomsMu.Lock()
	defer ss.roomsMu.Unlock()

	if len(ss.rooms) >= ss.config.MaxRooms {
		return fmt.Errorf("maximum number of rooms reached")
	}

	if _, exists := ss.rooms[info.ID]; exists {
		return fmt.Errorf("room %s already exists", info.ID)
	}

	room := newRoom(info)
	ss.rooms[info.ID] = room

	ss.statsMu.Lock()
	ss.stats.TotalRooms = len(ss.rooms)
	ss.statsMu.Unlock()

	ss.handlersMu.RLock()
	if ss.onRoomCreated != nil {
		go ss.onRoomCreated(info)
	}
	ss.handlersMu.RUnlock()

	return nil
}

func (ss *signalingServer) DeleteRoom(roomID string) error {
	ss.roomsMu.Lock()
	room, exists := ss.rooms[roomID]
	if exists {
		delete(ss.rooms, roomID)
		ss.statsMu.Lock()
		ss.stats.TotalRooms = len(ss.rooms)
		ss.statsMu.Unlock()
	}
	ss.roomsMu.Unlock()

	if exists {
		room.Close()

		ss.handlersMu.RLock()
		if ss.onRoomDeleted != nil {
			go ss.onRoomDeleted(roomID)
		}
		ss.handlersMu.RUnlock()
	}

	return nil
}

func (ss *signalingServer) GetRoom(roomID string) (*RoomInfo, error) {
	ss.roomsMu.RLock()
	room, exists := ss.rooms[roomID]
	ss.roomsMu.RUnlock()

	if !exists {
		return nil, ErrRoomNotFound
	}

	info := room.Info()
	return info, nil
}

func (ss *signalingServer) ListRooms() ([]*RoomInfo, error) {
	ss.roomsMu.RLock()
	defer ss.roomsMu.RUnlock()

	rooms := make([]*RoomInfo, 0, len(ss.rooms))
	for _, room := range ss.rooms {
		rooms = append(rooms, room.Info())
	}

	return rooms, nil
}

func (ss *signalingServer) GetRoomPeers(roomID string) ([]*PeerInfo, error) {
	ss.roomsMu.RLock()
	room, exists := ss.rooms[roomID]
	ss.roomsMu.RUnlock()

	if !exists {
		return nil, ErrRoomNotFound
	}

	return room.GetPeers(), nil
}

// Peer management
func (ss *signalingServer) GetPeer(peerID string) (*PeerInfo, error) {
	ss.peersMu.RLock()
	peer, exists := ss.peers[peerID]
	ss.peersMu.RUnlock()

	if !exists {
		return nil, ErrPeerNotFound
	}

	return peer.info, nil
}

func (ss *signalingServer) DisconnectPeer(peerID string) error {
	ss.peersMu.RLock()
	peer, exists := ss.peers[peerID]
	ss.peersMu.RUnlock()

	if !exists {
		return ErrPeerNotFound
	}

	peer.cancel()
	peer.conn.Close()

	return nil
}

func (ss *signalingServer) BroadcastToRoom(roomID string, msg *SignalingMessage) error {
	ss.roomsMu.RLock()
	room, exists := ss.rooms[roomID]
	ss.roomsMu.RUnlock()

	if !exists {
		return ErrRoomNotFound
	}

	return room.Broadcast(msg)
}

func (ss *signalingServer) SendToPeer(peerID string, msg *SignalingMessage) error {
	ss.peersMu.RLock()
	peer, exists := ss.peers[peerID]
	ss.peersMu.RUnlock()

	if !exists {
		return ErrPeerNotFound
	}

	select {
	case peer.send <- msg:
		return nil
	default:
		return fmt.Errorf("peer send queue is full")
	}
}

// Event handlers
func (ss *signalingServer) OnPeerConnected(handler func(*PeerInfo)) {
	ss.handlersMu.Lock()
	ss.onPeerConnected = handler
	ss.handlersMu.Unlock()
}

func (ss *signalingServer) OnPeerDisconnected(handler func(string)) {
	ss.handlersMu.Lock()
	ss.onPeerDisconnected = handler
	ss.handlersMu.Unlock()
}

func (ss *signalingServer) OnRoomCreated(handler func(*RoomInfo)) {
	ss.handlersMu.Lock()
	ss.onRoomCreated = handler
	ss.handlersMu.Unlock()
}

func (ss *signalingServer) OnRoomDeleted(handler func(string)) {
	ss.handlersMu.Lock()
	ss.onRoomDeleted = handler
	ss.handlersMu.Unlock()
}

func (ss *signalingServer) OnMessage(handler func(*SignalingMessage)) {
	ss.handlersMu.Lock()
	ss.onMessage = handler
	ss.handlersMu.Unlock()
}

func (ss *signalingServer) OnError(handler func(error)) {
	ss.handlersMu.Lock()
	ss.onError = handler
	ss.handlersMu.Unlock()
}

// Authentication
func (ss *signalingServer) SetAuthHandler(handler func(string) (*PeerInfo, error)) {
	ss.authHandler = handler
}

// Middleware
func (ss *signalingServer) AddMiddleware(middleware SignalingMiddleware) {
	ss.middlewares = append(ss.middlewares, middleware)
}

// Statistics
func (ss *signalingServer) GetStats() *ServerStats {
	ss.statsMu.RLock()
	defer ss.statsMu.RUnlock()

	stats := *ss.stats
	stats.Uptime = time.Since(ss.startTime).Milliseconds()
	stats.LastUpdated = time.Now()

	return &stats
}

// Configuration
func (ss *signalingServer) SetConfig(config *ServerConfig) {
	ss.config = config
}

func (ss *signalingServer) GetConfig() *ServerConfig {
	return ss.config
}

// statsCollector thu thập statistics định kỳ
func (ss *signalingServer) statsCollector() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ss.ctx.Done():
			return
		case <-ticker.C:
			ss.updateStats()
		}
	}
}

// updateStats cập nhật statistics
func (ss *signalingServer) updateStats() {
	ss.statsMu.Lock()
	defer ss.statsMu.Unlock()

	ss.peersMu.RLock()
	ss.stats.TotalPeers = len(ss.peers)
	ss.peersMu.RUnlock()

	ss.roomsMu.RLock()
	ss.stats.TotalRooms = len(ss.rooms)
	ss.roomsMu.RUnlock()
}

// emitError emit error event
func (ss *signalingServer) emitError(err error) {
	ss.handlersMu.RLock()
	if ss.onError != nil {
		go ss.onError(err)
	}
	ss.handlersMu.RUnlock()
}

// Peer methods
func (sp *signalingPeer) readPump() {
	defer func() {
		sp.server.removePeer(sp.info.ID)
		sp.conn.Close()
	}()

	sp.conn.SetReadLimit(512)
	sp.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	sp.conn.SetPongHandler(func(string) error {
		sp.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-sp.ctx.Done():
			return
		default:
		}

		_, data, err := sp.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sp.server.emitError(fmt.Errorf("websocket read error: %w", err))
			}
			return
		}

		var msg SignalingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			sp.server.emitError(fmt.Errorf("failed to unmarshal message: %w", err))
			continue
		}

		msg.From = sp.info.ID
		msg.Timestamp = time.Now()

		sp.handleMessage(&msg)
	}
}

func (sp *signalingPeer) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		sp.conn.Close()
	}()

	for {
		select {
		case <-sp.ctx.Done():
			return

		case msg := <-sp.send:
			sp.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			data, err := json.Marshal(msg)
			if err != nil {
				sp.server.emitError(fmt.Errorf("failed to marshal message: %w", err))
				continue
			}

			if err := sp.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				sp.server.emitError(fmt.Errorf("websocket write error: %w", err))
				return
			}

		case <-ticker.C:
			sp.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := sp.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (sp *signalingPeer) handleMessage(msg *SignalingMessage) {
	// Apply middleware
	handler := sp.processMessage
	for i := len(sp.server.middlewares) - 1; i >= 0; i-- {
		handler = sp.server.middlewares[i](handler)
	}

	if err := handler(msg); err != nil {
		sp.server.emitError(fmt.Errorf("message handling error: %w", err))
	}
}

func (sp *signalingPeer) processMessage(msg *SignalingMessage) error {
	// Emit message event
	sp.server.handlersMu.RLock()
	if sp.server.onMessage != nil {
		go sp.server.onMessage(msg)
	}
	sp.server.handlersMu.RUnlock()

	// Handle specific message types
	switch msg.Type {
	case "join-room":
		return sp.handleJoinRoom(msg)
	case "leave-room":
		return sp.handleLeaveRoom(msg)
	case "offer", "answer", "ice-candidate":
		return sp.handleSignalingMessage(msg)
	}

	return nil
}

func (sp *signalingPeer) handleJoinRoom(msg *SignalingMessage) error {
	if msg.Room == "" {
		return fmt.Errorf("room ID is required")
	}

	sp.server.roomsMu.RLock()
	room, exists := sp.server.rooms[msg.Room]
	sp.server.roomsMu.RUnlock()

	if !exists {
		// Create room if it doesn't exist
		roomInfo := &RoomInfo{
			ID:        msg.Room,
			Name:      msg.Room,
			MaxPeers:  sp.server.config.MaxPeersPerRoom,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := sp.server.CreateRoom(roomInfo); err != nil {
			return err
		}

		sp.server.roomsMu.RLock()
		room = sp.server.rooms[msg.Room]
		sp.server.roomsMu.RUnlock()
	}

	return room.AddPeer(sp.info)
}

func (sp *signalingPeer) handleLeaveRoom(msg *SignalingMessage) error {
	if msg.Room == "" {
		return fmt.Errorf("room ID is required")
	}

	sp.server.roomsMu.RLock()
	room, exists := sp.server.rooms[msg.Room]
	sp.server.roomsMu.RUnlock()

	if !exists {
		return ErrRoomNotFound
	}

	return room.RemovePeer(sp.info.ID)
}

func (sp *signalingPeer) handleSignalingMessage(msg *SignalingMessage) error {
	if msg.To == "" {
		return fmt.Errorf("target peer ID is required")
	}

	return sp.server.SendToPeer(msg.To, msg)
}

// removePeer xóa peer khỏi server
func (ss *signalingServer) removePeer(peerID string) {
	ss.peersMu.Lock()
	_, exists := ss.peers[peerID]
	if exists {
		delete(ss.peers, peerID)
		ss.statsMu.Lock()
		ss.stats.ActiveConnections--
		ss.statsMu.Unlock()
	}
	ss.peersMu.Unlock()

	if exists {
		// Remove peer from all rooms
		ss.roomsMu.RLock()
		for _, room := range ss.rooms {
			room.RemovePeer(peerID)
		}
		ss.roomsMu.RUnlock()

		ss.handlersMu.RLock()
		if ss.onPeerDisconnected != nil {
			go ss.onPeerDisconnected(peerID)
		}
		ss.handlersMu.RUnlock()
	}
}

// Helper functions
func generatePeerID() string {
	return fmt.Sprintf("peer_%d", time.Now().UnixNano())
}
