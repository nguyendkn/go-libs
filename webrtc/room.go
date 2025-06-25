package webrtc

import (
	"fmt"
	"sync"
	"time"
)

// room implements the Room interface
type room struct {
	info *RoomInfo

	// Peer management
	peers   map[string]*PeerInfo
	peersMu sync.RWMutex

	// Event handlers
	onPeerJoined func(*PeerInfo)
	onPeerLeft   func(string)
	onMessage    func(*SignalingMessage)
	onEmpty      func()
	handlersMu   sync.RWMutex

	// State
	active bool
	mu     sync.RWMutex
}

// newRoom tạo một Room mới
func newRoom(info *RoomInfo) Room {
	return &room{
		info:   info,
		peers:  make(map[string]*PeerInfo),
		active: true,
	}
}

// Room info
func (r *room) ID() string {
	return r.info.ID
}

func (r *room) Info() *RoomInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy of room info
	info := *r.info
	info.UpdatedAt = time.Now()

	return &info
}

// Peer management
func (r *room) AddPeer(peer *PeerInfo) error {
	r.peersMu.Lock()
	defer r.peersMu.Unlock()

	// Check if room is full
	if r.info.MaxPeers > 0 && len(r.peers) >= r.info.MaxPeers {
		return ErrRoomFull
	}

	// Check if peer already exists
	if _, exists := r.peers[peer.ID]; exists {
		return fmt.Errorf("peer %s already in room", peer.ID)
	}

	// Add peer
	r.peers[peer.ID] = peer

	// Update room info
	r.mu.Lock()
	r.info.UpdatedAt = time.Now()
	r.mu.Unlock()

	// Emit peer joined event
	r.handlersMu.RLock()
	if r.onPeerJoined != nil {
		go r.onPeerJoined(peer)
	}
	r.handlersMu.RUnlock()

	return nil
}

func (r *room) RemovePeer(peerID string) error {
	r.peersMu.Lock()
	_, exists := r.peers[peerID]
	if exists {
		delete(r.peers, peerID)
	}
	isEmpty := len(r.peers) == 0
	r.peersMu.Unlock()

	if !exists {
		return ErrPeerNotFound
	}

	// Update room info
	r.mu.Lock()
	r.info.UpdatedAt = time.Now()
	r.mu.Unlock()

	// Emit peer left event
	r.handlersMu.RLock()
	if r.onPeerLeft != nil {
		go r.onPeerLeft(peerID)
	}

	// Emit empty event if room is now empty
	if isEmpty && r.onEmpty != nil {
		go r.onEmpty()
	}
	r.handlersMu.RUnlock()

	return nil
}

func (r *room) GetPeer(peerID string) (*PeerInfo, error) {
	r.peersMu.RLock()
	peer, exists := r.peers[peerID]
	r.peersMu.RUnlock()

	if !exists {
		return nil, ErrPeerNotFound
	}

	return peer, nil
}

func (r *room) GetPeers() []*PeerInfo {
	r.peersMu.RLock()
	defer r.peersMu.RUnlock()

	peers := make([]*PeerInfo, 0, len(r.peers))
	for _, peer := range r.peers {
		peers = append(peers, peer)
	}

	return peers
}

func (r *room) GetPeerCount() int {
	r.peersMu.RLock()
	defer r.peersMu.RUnlock()

	return len(r.peers)
}

func (r *room) HasPeer(peerID string) bool {
	r.peersMu.RLock()
	_, exists := r.peers[peerID]
	r.peersMu.RUnlock()

	return exists
}

// Broadcasting
func (r *room) Broadcast(msg *SignalingMessage) error {
	r.peersMu.RLock()
	defer r.peersMu.RUnlock()

	// This is a simplified implementation
	// In a real implementation, you would need access to the signaling server
	// to actually send messages to peers

	// Emit message event
	r.handlersMu.RLock()
	if r.onMessage != nil {
		go r.onMessage(msg)
	}
	r.handlersMu.RUnlock()

	return nil
}

func (r *room) BroadcastExcept(msg *SignalingMessage, excludePeerID string) error {
	// Similar to Broadcast but exclude specific peer
	return r.Broadcast(msg)
}

func (r *room) SendToPeer(peerID string, msg *SignalingMessage) error {
	r.peersMu.RLock()
	_, exists := r.peers[peerID]
	r.peersMu.RUnlock()

	if !exists {
		return ErrPeerNotFound
	}

	// This is a simplified implementation
	// In a real implementation, you would need access to the signaling server
	// to actually send the message to the peer

	return nil
}

// Room state
func (r *room) IsEmpty() bool {
	return r.GetPeerCount() == 0
}

func (r *room) IsFull() bool {
	if r.info.MaxPeers <= 0 {
		return false // No limit
	}
	return r.GetPeerCount() >= r.info.MaxPeers
}

func (r *room) IsActive() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.active
}

// Event handlers
func (r *room) OnPeerJoined(handler func(*PeerInfo)) {
	r.handlersMu.Lock()
	r.onPeerJoined = handler
	r.handlersMu.Unlock()
}

func (r *room) OnPeerLeft(handler func(string)) {
	r.handlersMu.Lock()
	r.onPeerLeft = handler
	r.handlersMu.Unlock()
}

func (r *room) OnMessage(handler func(*SignalingMessage)) {
	r.handlersMu.Lock()
	r.onMessage = handler
	r.handlersMu.Unlock()
}

func (r *room) OnEmpty(handler func()) {
	r.handlersMu.Lock()
	r.onEmpty = handler
	r.handlersMu.Unlock()
}

// Configuration
func (r *room) UpdateInfo(info *RoomInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Preserve ID and creation time
	info.ID = r.info.ID
	info.CreatedAt = r.info.CreatedAt
	info.UpdatedAt = time.Now()

	r.info = info
	return nil
}

func (r *room) SetMaxPeers(max int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if current peer count exceeds new limit
	currentCount := r.GetPeerCount()
	if max > 0 && currentCount > max {
		return fmt.Errorf("current peer count (%d) exceeds new limit (%d)", currentCount, max)
	}

	r.info.MaxPeers = max
	r.info.UpdatedAt = time.Now()

	return nil
}

// Lifecycle
func (r *room) Close() error {
	r.mu.Lock()
	r.active = false
	r.mu.Unlock()

	// Remove all peers
	r.peersMu.Lock()
	peerIDs := make([]string, 0, len(r.peers))
	for peerID := range r.peers {
		peerIDs = append(peerIDs, peerID)
	}
	r.peersMu.Unlock()

	// Remove peers one by one to trigger events
	for _, peerID := range peerIDs {
		r.RemovePeer(peerID)
	}

	return nil
}
