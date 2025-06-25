package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// wsRoom implements the Room interface
type wsRoom struct {
	name    string
	options *RoomOptions

	// Client management
	clients   map[string]ServerClient
	clientsMu sync.RWMutex

	// Message history
	history   []*Message
	historyMu sync.RWMutex

	// Event handlers
	onClientJoin  func(ServerClient)
	onClientLeave func(ServerClient)
	onMessage     func(ServerClient, []byte)
	onEmpty       func()
	handlersMu    sync.RWMutex

	// State
	closed   bool
	closedMu sync.RWMutex

	// Creation time
	createdAt time.Time
}

// newRoom tạo một room mới
func newRoom(name string, options *RoomOptions) Room {
	if options == nil {
		options = &RoomOptions{
			MaxClients:     0, // unlimited
			RequireAuth:    false,
			MessageHistory: 100,
			TTL:            0, // no TTL
		}
	}

	return &wsRoom{
		name:      name,
		options:   options,
		clients:   make(map[string]ServerClient),
		history:   make([]*Message, 0, options.MessageHistory),
		createdAt: time.Now(),
	}
}

// Basic info methods
func (r *wsRoom) Name() string {
	return r.name
}

func (r *wsRoom) Options() *RoomOptions {
	return r.options
}

// Client management methods
func (r *wsRoom) AddClient(client ServerClient) error {
	r.closedMu.RLock()
	if r.closed {
		r.closedMu.RUnlock()
		return fmt.Errorf("room is closed")
	}
	r.closedMu.RUnlock()

	// Check if room is full
	if r.options.MaxClients > 0 && r.GetClientCount() >= r.options.MaxClients {
		return ErrRoomFull
	}

	// Check authentication if required
	if r.options.RequireAuth && !client.IsAuthenticated() {
		return ErrUnauthorized
	}

	// Check roles if specified
	if len(r.options.AllowedRoles) > 0 {
		auth := client.GetAuth()
		if auth == nil {
			return ErrUnauthorized
		}

		hasRole := false
		for _, allowedRole := range r.options.AllowedRoles {
			for _, userRole := range auth.Roles {
				if userRole == allowedRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return ErrUnauthorized
		}
	}

	r.clientsMu.Lock()
	r.clients[client.ID()] = client
	r.clientsMu.Unlock()

	// Send message history to new client
	if r.options.MessageHistory > 0 {
		r.historyMu.RLock()
		for _, msg := range r.history {
			client.SendWithType(msg.Type, msg.Data)
		}
		r.historyMu.RUnlock()
	}

	// Call join handler
	r.handlersMu.RLock()
	if r.onClientJoin != nil {
		go r.onClientJoin(client)
	}
	r.handlersMu.RUnlock()

	return nil
}

func (r *wsRoom) RemoveClient(clientID string) error {
	r.clientsMu.Lock()
	client, exists := r.clients[clientID]
	if exists {
		delete(r.clients, clientID)
	}
	isEmpty := len(r.clients) == 0
	r.clientsMu.Unlock()

	if exists {
		// Call leave handler
		r.handlersMu.RLock()
		if r.onClientLeave != nil {
			go r.onClientLeave(client)
		}
		if isEmpty && r.onEmpty != nil {
			go r.onEmpty()
		}
		r.handlersMu.RUnlock()
	}

	return nil
}

func (r *wsRoom) GetClient(clientID string) (ServerClient, bool) {
	r.clientsMu.RLock()
	client, exists := r.clients[clientID]
	r.clientsMu.RUnlock()
	return client, exists
}

func (r *wsRoom) GetClients() []ServerClient {
	r.clientsMu.RLock()
	clients := make([]ServerClient, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	r.clientsMu.RUnlock()
	return clients
}

func (r *wsRoom) GetClientCount() int {
	r.clientsMu.RLock()
	count := len(r.clients)
	r.clientsMu.RUnlock()
	return count
}

func (r *wsRoom) HasClient(clientID string) bool {
	r.clientsMu.RLock()
	_, exists := r.clients[clientID]
	r.clientsMu.RUnlock()
	return exists
}

// Broadcasting methods
func (r *wsRoom) Broadcast(data []byte) error {
	return r.BroadcastWithType(BinaryMessage, data)
}

func (r *wsRoom) BroadcastText(text string) error {
	return r.BroadcastWithType(TextMessage, []byte(text))
}

func (r *wsRoom) BroadcastJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return r.BroadcastText(string(data))
}

func (r *wsRoom) BroadcastExcept(data []byte, excludeClientID string) error {
	return r.BroadcastExceptWithType(BinaryMessage, data, excludeClientID)
}

// BroadcastWithType broadcasts message with specific type
func (r *wsRoom) BroadcastWithType(messageType MessageType, data []byte) error {
	r.closedMu.RLock()
	if r.closed {
		r.closedMu.RUnlock()
		return fmt.Errorf("room is closed")
	}
	r.closedMu.RUnlock()

	// Create message
	msg := &Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		Room:      r.name,
	}

	// Add to history
	r.AddToHistory(msg)

	// Broadcast to all clients
	r.clientsMu.RLock()
	clients := make([]ServerClient, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	r.clientsMu.RUnlock()

	for _, client := range clients {
		if client.IsConnected() {
			client.SendWithType(messageType, data)
		}
	}

	return nil
}

// BroadcastExceptWithType broadcasts message except to specific client
func (r *wsRoom) BroadcastExceptWithType(messageType MessageType, data []byte, excludeClientID string) error {
	r.closedMu.RLock()
	if r.closed {
		r.closedMu.RUnlock()
		return fmt.Errorf("room is closed")
	}
	r.closedMu.RUnlock()

	// Create message
	msg := &Message{
		Type:      messageType,
		Data:      data,
		Timestamp: time.Now(),
		Room:      r.name,
	}

	// Add to history
	r.AddToHistory(msg)

	// Broadcast to all clients except excluded one
	r.clientsMu.RLock()
	clients := make([]ServerClient, 0, len(r.clients))
	for _, client := range r.clients {
		if client.ID() != excludeClientID {
			clients = append(clients, client)
		}
	}
	r.clientsMu.RUnlock()

	for _, client := range clients {
		if client.IsConnected() {
			client.SendWithType(messageType, data)
		}
	}

	return nil
}

// Message history methods
func (r *wsRoom) GetMessageHistory() []*Message {
	r.historyMu.RLock()
	history := make([]*Message, len(r.history))
	copy(history, r.history)
	r.historyMu.RUnlock()
	return history
}

func (r *wsRoom) AddToHistory(msg *Message) {
	if r.options.MessageHistory <= 0 {
		return
	}

	r.historyMu.Lock()
	r.history = append(r.history, msg)

	// Keep only the last N messages
	if len(r.history) > r.options.MessageHistory {
		r.history = r.history[len(r.history)-r.options.MessageHistory:]
	}
	r.historyMu.Unlock()
}

func (r *wsRoom) ClearHistory() {
	r.historyMu.Lock()
	r.history = r.history[:0]
	r.historyMu.Unlock()
}

// Room state methods
func (r *wsRoom) IsEmpty() bool {
	return r.GetClientCount() == 0
}

func (r *wsRoom) IsFull() bool {
	if r.options.MaxClients <= 0 {
		return false // unlimited
	}
	return r.GetClientCount() >= r.options.MaxClients
}

// Lifecycle methods
func (r *wsRoom) Close() error {
	r.closedMu.Lock()
	if r.closed {
		r.closedMu.Unlock()
		return nil
	}
	r.closed = true
	r.closedMu.Unlock()

	// Disconnect all clients
	r.clientsMu.Lock()
	clients := make([]ServerClient, 0, len(r.clients))
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	r.clients = make(map[string]ServerClient) // clear the map
	r.clientsMu.Unlock()

	// Close all client connections
	for _, client := range clients {
		client.CloseWithCode(websocket.CloseNormalClosure, "room closed")
	}

	// Clear history
	r.ClearHistory()

	return nil
}

// Event handlers
func (r *wsRoom) OnClientJoin(handler func(ServerClient)) {
	r.handlersMu.Lock()
	r.onClientJoin = handler
	r.handlersMu.Unlock()
}

func (r *wsRoom) OnClientLeave(handler func(ServerClient)) {
	r.handlersMu.Lock()
	r.onClientLeave = handler
	r.handlersMu.Unlock()
}

func (r *wsRoom) OnMessage(handler func(ServerClient, []byte)) {
	r.handlersMu.Lock()
	r.onMessage = handler
	r.handlersMu.Unlock()
}

func (r *wsRoom) OnEmpty(handler func()) {
	r.handlersMu.Lock()
	r.onEmpty = handler
	r.handlersMu.Unlock()
}
