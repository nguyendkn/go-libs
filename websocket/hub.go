package websocket

import (
	"fmt"
	"sync"
	"time"
)

// wsHub implements the Hub interface
type wsHub struct {
	// Client management
	clients   map[string]ServerClient
	clientsMu sync.RWMutex

	// Room management
	rooms   map[string]Room
	roomsMu sync.RWMutex

	// Channels
	register   chan ServerClient
	unregister chan string
	broadcast  chan *Message
	roomcast   chan *roomMessage

	// State
	running bool
	runMu   sync.RWMutex

	// Metrics
	metrics *HubMetrics

	// Wait group
	wg sync.WaitGroup
}

type roomMessage struct {
	room string
	msg  *Message
}

// NewHub tạo một hub mới
func NewHub() Hub {
	return &wsHub{
		clients:    make(map[string]ServerClient),
		rooms:      make(map[string]Room),
		register:   make(chan ServerClient, 100),
		unregister: make(chan string, 100),
		broadcast:  make(chan *Message, 1000),
		roomcast:   make(chan *roomMessage, 1000),
		metrics: &HubMetrics{
			RoomMetrics: make(map[string]*RoomMetrics),
		},
	}
}

// Start khởi động hub
func (h *wsHub) Start() error {
	h.runMu.Lock()
	if h.running {
		h.runMu.Unlock()
		return fmt.Errorf("hub is already running")
	}
	h.running = true
	h.runMu.Unlock()

	h.wg.Add(1)
	go h.run()

	return nil
}

// Stop dừng hub
func (h *wsHub) Stop() error {
	h.runMu.Lock()
	if !h.running {
		h.runMu.Unlock()
		return nil
	}
	h.running = false
	h.runMu.Unlock()

	// Close channels
	close(h.register)
	close(h.unregister)
	close(h.broadcast)
	close(h.roomcast)

	// Wait for goroutines to finish
	h.wg.Wait()

	// Close all rooms
	h.roomsMu.Lock()
	for _, room := range h.rooms {
		room.Close()
	}
	h.roomsMu.Unlock()

	return nil
}

// run chạy hub loop chính
func (h *wsHub) run() {
	defer h.wg.Done()

	for {
		select {
		case client, ok := <-h.register:
			if !ok {
				return
			}
			h.registerClient(client)

		case clientID, ok := <-h.unregister:
			if !ok {
				return
			}
			h.unregisterClient(clientID)

		case msg, ok := <-h.broadcast:
			if !ok {
				return
			}
			h.broadcastMessage(msg)

		case roomMsg, ok := <-h.roomcast:
			if !ok {
				return
			}
			h.broadcastToRoom(roomMsg.room, roomMsg.msg)
		}
	}
}

// registerClient đăng ký client
func (h *wsHub) registerClient(client ServerClient) {
	h.clientsMu.Lock()
	h.clients[client.ID()] = client
	h.metrics.RegisteredClients = int64(len(h.clients))
	h.clientsMu.Unlock()
}

// unregisterClient hủy đăng ký client
func (h *wsHub) unregisterClient(clientID string) {
	h.clientsMu.Lock()
	client, exists := h.clients[clientID]
	if exists {
		delete(h.clients, clientID)
		h.metrics.RegisteredClients = int64(len(h.clients))
	}
	h.clientsMu.Unlock()

	if exists {
		// Remove client from all rooms
		rooms := client.GetRooms()
		for _, roomName := range rooms {
			h.roomsMu.RLock()
			if room, ok := h.rooms[roomName]; ok {
				room.RemoveClient(clientID)
			}
			h.roomsMu.RUnlock()
		}

		client.Close()
	}
}

// broadcastMessage broadcast message to all clients
func (h *wsHub) broadcastMessage(msg *Message) {
	h.clientsMu.RLock()
	clients := make([]ServerClient, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	h.clientsMu.RUnlock()

	for _, client := range clients {
		if client.IsConnected() {
			client.SendWithType(msg.Type, msg.Data)
		}
	}

	h.metrics.MessagesProcessed++
}

// broadcastToRoom broadcast message to room
func (h *wsHub) broadcastToRoom(roomName string, msg *Message) {
	h.roomsMu.RLock()
	room, exists := h.rooms[roomName]
	h.roomsMu.RUnlock()

	if exists {
		room.Broadcast(msg.Data)
		h.metrics.MessagesProcessed++
	}
}

// Client management methods
func (h *wsHub) RegisterClient(client ServerClient) {
	select {
	case h.register <- client:
	default:
		// Channel is full, register directly
		h.registerClient(client)
	}
}

func (h *wsHub) UnregisterClient(clientID string) {
	select {
	case h.unregister <- clientID:
	default:
		// Channel is full, unregister directly
		h.unregisterClient(clientID)
	}
}

func (h *wsHub) GetClient(clientID string) (ServerClient, bool) {
	h.clientsMu.RLock()
	client, exists := h.clients[clientID]
	h.clientsMu.RUnlock()
	return client, exists
}

func (h *wsHub) GetClients() []ServerClient {
	h.clientsMu.RLock()
	clients := make([]ServerClient, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	h.clientsMu.RUnlock()
	return clients
}

func (h *wsHub) GetClientCount() int {
	h.clientsMu.RLock()
	count := len(h.clients)
	h.clientsMu.RUnlock()
	return count
}

// Room management methods
func (h *wsHub) CreateRoom(name string, options *RoomOptions) error {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	if _, exists := h.rooms[name]; exists {
		return fmt.Errorf("room %s already exists", name)
	}

	room := newRoom(name, options)
	h.rooms[name] = room
	h.metrics.ActiveRooms = len(h.rooms)
	h.metrics.RoomMetrics[name] = &RoomMetrics{
		Name:         name,
		ClientCount:  0,
		MessageCount: 0,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	return nil
}

func (h *wsHub) DeleteRoom(name string) error {
	h.roomsMu.Lock()
	room, exists := h.rooms[name]
	if exists {
		delete(h.rooms, name)
		delete(h.metrics.RoomMetrics, name)
		h.metrics.ActiveRooms = len(h.rooms)
	}
	h.roomsMu.Unlock()

	if exists {
		room.Close()
	}

	return nil
}

func (h *wsHub) GetRoom(name string) (Room, bool) {
	h.roomsMu.RLock()
	room, exists := h.rooms[name]
	h.roomsMu.RUnlock()
	return room, exists
}

func (h *wsHub) GetRooms() []Room {
	h.roomsMu.RLock()
	rooms := make([]Room, 0, len(h.rooms))
	for _, room := range h.rooms {
		rooms = append(rooms, room)
	}
	h.roomsMu.RUnlock()
	return rooms
}

func (h *wsHub) JoinRoom(clientID, roomName string) error {
	// Get client
	client, exists := h.GetClient(clientID)
	if !exists {
		return fmt.Errorf("client %s not found", clientID)
	}

	// Get or create room
	h.roomsMu.RLock()
	room, exists := h.rooms[roomName]
	h.roomsMu.RUnlock()

	if !exists {
		// Create room with default options
		if err := h.CreateRoom(roomName, &RoomOptions{}); err != nil {
			return err
		}
		h.roomsMu.RLock()
		room = h.rooms[roomName]
		h.roomsMu.RUnlock()
	}

	// Add client to room
	if err := room.AddClient(client); err != nil {
		return err
	}

	// Add room to client
	return client.JoinRoom(roomName)
}

func (h *wsHub) LeaveRoom(clientID, roomName string) error {
	// Get client
	client, exists := h.GetClient(clientID)
	if !exists {
		return fmt.Errorf("client %s not found", clientID)
	}

	// Get room
	h.roomsMu.RLock()
	room, exists := h.rooms[roomName]
	h.roomsMu.RUnlock()

	if !exists {
		return fmt.Errorf("room %s not found", roomName)
	}

	// Remove client from room
	room.RemoveClient(clientID)

	// Remove room from client
	return client.LeaveRoom(roomName)
}

// Broadcasting methods
func (h *wsHub) Broadcast(msg *Message) error {
	select {
	case h.broadcast <- msg:
		return nil
	default:
		// Channel is full, broadcast directly
		h.broadcastMessage(msg)
		return nil
	}
}

func (h *wsHub) BroadcastToRoom(roomName string, msg *Message) error {
	roomMsg := &roomMessage{
		room: roomName,
		msg:  msg,
	}

	select {
	case h.roomcast <- roomMsg:
		return nil
	default:
		// Channel is full, broadcast directly
		h.broadcastToRoom(roomName, msg)
		return nil
	}
}

func (h *wsHub) BroadcastToClient(clientID string, msg *Message) error {
	client, exists := h.GetClient(clientID)
	if !exists {
		return fmt.Errorf("client %s not found", clientID)
	}

	return client.SendWithType(msg.Type, msg.Data)
}

// Event handling
func (h *wsHub) HandleEvent(event *Event) {
	h.metrics.EventsProcessed++
	// Event handling logic can be implemented here
}

// Monitoring
func (h *wsHub) GetMetrics() *HubMetrics {
	h.clientsMu.RLock()
	h.roomsMu.RLock()

	metrics := &HubMetrics{
		RegisteredClients: int64(len(h.clients)),
		ActiveRooms:       len(h.rooms),
		MessagesProcessed: h.metrics.MessagesProcessed,
		EventsProcessed:   h.metrics.EventsProcessed,
		QueueSize:         len(h.broadcast) + len(h.roomcast),
		RoomMetrics:       make(map[string]*RoomMetrics),
	}

	// Copy room metrics
	for name, roomMetrics := range h.metrics.RoomMetrics {
		metrics.RoomMetrics[name] = &RoomMetrics{
			Name:         roomMetrics.Name,
			ClientCount:  roomMetrics.ClientCount,
			MessageCount: roomMetrics.MessageCount,
			CreatedAt:    roomMetrics.CreatedAt,
			LastActivity: roomMetrics.LastActivity,
		}
	}

	h.roomsMu.RUnlock()
	h.clientsMu.RUnlock()

	return metrics
}
