package webrtc

import (
	"fmt"
	"time"
)

// ConnectionState định nghĩa trạng thái kết nối WebRTC
type ConnectionState int

const (
	ConnectionStateNew ConnectionState = iota
	ConnectionStateConnecting
	ConnectionStateConnected
	ConnectionStateDisconnected
	ConnectionStateFailed
	ConnectionStateClosed
)

// ICEConnectionState định nghĩa trạng thái ICE connection
type ICEConnectionState int

const (
	ICEConnectionStateNew ICEConnectionState = iota
	ICEConnectionStateChecking
	ICEConnectionStateConnected
	ICEConnectionStateCompleted
	ICEConnectionStateDisconnected
	ICEConnectionStateFailed
	ICEConnectionStateClosed
)

// SignalingState định nghĩa trạng thái signaling
type SignalingState int

const (
	SignalingStateStable SignalingState = iota
	SignalingStateHaveLocalOffer
	SignalingStateHaveRemoteOffer
	SignalingStateHaveLocalPranswer
	SignalingStateHaveRemotePranswer
	SignalingStateClosed
)

// DataChannelState định nghĩa trạng thái DataChannel
type DataChannelState int

const (
	DataChannelStateConnecting DataChannelState = iota
	DataChannelStateOpen
	DataChannelStateClosing
	DataChannelStateClosed
)

// MediaType định nghĩa loại media
type MediaType int

const (
	MediaTypeAudio MediaType = iota
	MediaTypeVideo
	MediaTypeData
)

// TrackDirection định nghĩa hướng track
type TrackDirection int

const (
	TrackDirectionSendOnly TrackDirection = iota
	TrackDirectionRecvOnly
	TrackDirectionSendRecv
	TrackDirectionInactive
)

// SessionDescription đại diện cho SDP
type SessionDescription struct {
	Type string `json:"type"` // "offer", "answer", "pranswer", "rollback"
	SDP  string `json:"sdp"`
}

// ICECandidate đại diện cho ICE candidate
type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex uint16 `json:"sdpMLineIndex"`
}

// ICEServer cấu hình ICE server (STUN/TURN)
type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// PeerConnectionConfig cấu hình cho PeerConnection
type PeerConnectionConfig struct {
	ICEServers           []ICEServer `json:"iceServers"`
	ICETransportPolicy   string      `json:"iceTransportPolicy,omitempty"` // "all" or "relay"
	BundlePolicy         string      `json:"bundlePolicy,omitempty"`       // "balanced", "max-compat", "max-bundle"
	RTCPMuxPolicy        string      `json:"rtcpMuxPolicy,omitempty"`      // "negotiate" or "require"
	PeerIdentity         string      `json:"peerIdentity,omitempty"`
	Certificates         []string    `json:"certificates,omitempty"`
	ICECandidatePoolSize int         `json:"iceCandidatePoolSize,omitempty"`
	SDPSemantics         string      `json:"sdpSemantics,omitempty"` // "plan-b" or "unified-plan"

	// Custom options
	ConnectionTimeout   time.Duration `json:"connectionTimeout,omitempty"`
	DisconnectedTimeout time.Duration `json:"disconnectedTimeout,omitempty"`
	FailedTimeout       time.Duration `json:"failedTimeout,omitempty"`
	KeepAliveInterval   time.Duration `json:"keepAliveInterval,omitempty"`
}

// DataChannelConfig cấu hình cho DataChannel
type DataChannelConfig struct {
	Label             string `json:"label"`
	Protocol          string `json:"protocol,omitempty"`
	Negotiated        bool   `json:"negotiated,omitempty"`
	ID                uint16 `json:"id,omitempty"`
	Ordered           bool   `json:"ordered"`
	MaxPacketLifeTime uint16 `json:"maxPacketLifeTime,omitempty"`
	MaxRetransmits    uint16 `json:"maxRetransmits,omitempty"`
}

// MediaStreamTrack đại diện cho media track
type MediaStreamTrack struct {
	ID         string         `json:"id"`
	Kind       MediaType      `json:"kind"`
	Label      string         `json:"label"`
	Enabled    bool           `json:"enabled"`
	Muted      bool           `json:"muted"`
	ReadyState string         `json:"readyState"` // "live" or "ended"
	Direction  TrackDirection `json:"direction"`

	// Internal track reference
	TrackRef interface{} `json:"-"`
}

// MediaStream đại diện cho media stream
type MediaStream struct {
	ID     string              `json:"id"`
	Label  string              `json:"label"`
	Tracks []*MediaStreamTrack `json:"tracks"`
	Active bool                `json:"active"`
}

// RTCStatsReport đại diện cho WebRTC stats
type RTCStatsReport struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Stats     map[string]interface{} `json:"stats"`
}

// PeerConnectionStats thống kê kết nối
type PeerConnectionStats struct {
	ConnectionState    ConnectionState    `json:"connectionState"`
	ICEConnectionState ICEConnectionState `json:"iceConnectionState"`
	SignalingState     SignalingState     `json:"signalingState"`

	// Timing
	ConnectedAt  time.Time `json:"connectedAt"`
	LastActivity time.Time `json:"lastActivity"`

	// Data transfer
	BytesSent       uint64 `json:"bytesSent"`
	BytesReceived   uint64 `json:"bytesReceived"`
	PacketsSent     uint64 `json:"packetsSent"`
	PacketsReceived uint64 `json:"packetsReceived"`
	PacketsLost     uint64 `json:"packetsLost"`

	// Quality metrics
	RTT            time.Duration `json:"rtt"`            // Round Trip Time
	Jitter         time.Duration `json:"jitter"`         // Jitter
	PacketLossRate float64       `json:"packetLossRate"` // Packet loss rate (0-1)

	// Bandwidth
	AvailableOutgoingBitrate uint32 `json:"availableOutgoingBitrate"`
	AvailableIncomingBitrate uint32 `json:"availableIncomingBitrate"`

	// ICE
	LocalCandidates       []ICECandidate `json:"localCandidates"`
	RemoteCandidates      []ICECandidate `json:"remoteCandidates"`
	SelectedCandidatePair string         `json:"selectedCandidatePair"`
}

// SignalingMessage đại diện cho signaling message
type SignalingMessage struct {
	Type      string      `json:"type"` // "offer", "answer", "ice-candidate", "bye"
	From      string      `json:"from"`
	To        string      `json:"to"`
	Room      string      `json:"room,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// RoomInfo thông tin room
type RoomInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MaxPeers    int       `json:"maxPeers"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Room settings
	RequireAuth  bool     `json:"requireAuth"`
	AllowedRoles []string `json:"allowedRoles,omitempty"`
	IsPrivate    bool     `json:"isPrivate"`
	Password     string   `json:"password,omitempty"`

	// Media settings
	AudioEnabled bool `json:"audioEnabled"`
	VideoEnabled bool `json:"videoEnabled"`
	DataEnabled  bool `json:"dataEnabled"`

	// Quality settings
	MaxBitrate   uint32 `json:"maxBitrate,omitempty"`
	MaxFramerate uint32 `json:"maxFramerate,omitempty"`
	Resolution   string `json:"resolution,omitempty"` // "720p", "1080p", etc.
}

// PeerInfo thông tin peer trong room
type PeerInfo struct {
	ID       string    `json:"id"`
	UserID   string    `json:"userId,omitempty"`
	Username string    `json:"username,omitempty"`
	Role     string    `json:"role,omitempty"`
	JoinedAt time.Time `json:"joinedAt"`

	// Media state
	AudioEnabled bool `json:"audioEnabled"`
	VideoEnabled bool `json:"videoEnabled"`
	ScreenShare  bool `json:"screenShare"`

	// Connection info
	ConnectionState ConnectionState    `json:"connectionState"`
	ICEState        ICEConnectionState `json:"iceState"`

	// Quality metrics
	Quality struct {
		Audio struct {
			Bitrate    uint32  `json:"bitrate"`
			PacketLoss float64 `json:"packetLoss"`
		} `json:"audio"`
		Video struct {
			Bitrate    uint32  `json:"bitrate"`
			Framerate  uint32  `json:"framerate"`
			Resolution string  `json:"resolution"`
			PacketLoss float64 `json:"packetLoss"`
		} `json:"video"`
	} `json:"quality"`
}

// Event types
type EventType string

const (
	EventPeerConnected         EventType = "peer_connected"
	EventPeerDisconnected      EventType = "peer_disconnected"
	EventTrackAdded            EventType = "track_added"
	EventTrackRemoved          EventType = "track_removed"
	EventDataChannelOpen       EventType = "datachannel_open"
	EventDataChannelClose      EventType = "datachannel_close"
	EventDataChannelMessage    EventType = "datachannel_message"
	EventICECandidate          EventType = "ice_candidate"
	EventConnectionStateChange EventType = "connection_state_change"
	EventSignalingStateChange  EventType = "signaling_state_change"
	EventRoomJoined            EventType = "room_joined"
	EventRoomLeft              EventType = "room_left"
	EventError                 EventType = "error"
)

// Event đại diện cho WebRTC event
type Event struct {
	Type      EventType   `json:"type"`
	PeerID    string      `json:"peerId,omitempty"`
	RoomID    string      `json:"roomId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Error     error       `json:"error,omitempty"`
}

// Default values
const (
	DefaultConnectionTimeout   = 30 * time.Second
	DefaultDisconnectedTimeout = 5 * time.Second
	DefaultFailedTimeout       = 30 * time.Second
	DefaultKeepAliveInterval   = 25 * time.Second
	DefaultMaxBitrate          = 2500000 // 2.5 Mbps
	DefaultMaxFramerate        = 30
	DefaultMaxPeersPerRoom     = 50
	DefaultStatsInterval       = 5 * time.Second
)

// Common STUN/TURN servers
var (
	DefaultSTUNServers = []ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
		{URLs: []string{"stun:stun1.l.google.com:19302"}},
		{URLs: []string{"stun:stun2.l.google.com:19302"}},
	}

	DefaultICEServers = []ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
	}
)

// Error definitions
type WebRTCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *WebRTCError) Error() string {
	return fmt.Sprintf("WebRTC Error [%d]: %s", e.Code, e.Message)
}

// Common errors
var (
	ErrPeerConnectionClosed      = &WebRTCError{Code: 1001, Message: "peer connection is closed", Type: "connection"}
	ErrDataChannelClosed         = &WebRTCError{Code: 1002, Message: "data channel is closed", Type: "datachannel"}
	ErrInvalidSessionDescription = &WebRTCError{Code: 1003, Message: "invalid session description", Type: "sdp"}
	ErrInvalidICECandidate       = &WebRTCError{Code: 1004, Message: "invalid ICE candidate", Type: "ice"}
	ErrRoomNotFound              = &WebRTCError{Code: 1005, Message: "room not found", Type: "room"}
	ErrRoomFull                  = &WebRTCError{Code: 1006, Message: "room is full", Type: "room"}
	ErrPeerNotFound              = &WebRTCError{Code: 1007, Message: "peer not found", Type: "peer"}
	ErrUnauthorized              = &WebRTCError{Code: 1008, Message: "unauthorized", Type: "auth"}
	ErrMediaNotSupported         = &WebRTCError{Code: 1009, Message: "media type not supported", Type: "media"}
	ErrSignalingFailed           = &WebRTCError{Code: 1010, Message: "signaling failed", Type: "signaling"}
)
