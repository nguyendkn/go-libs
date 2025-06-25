package webrtc

import (
	"context"
	"io"
	"time"
)

// PeerConnection interface định nghĩa WebRTC peer connection
type PeerConnection interface {
	// Connection lifecycle
	CreateOffer(options *OfferOptions) (*SessionDescription, error)
	CreateAnswer(options *AnswerOptions) (*SessionDescription, error)
	SetLocalDescription(desc *SessionDescription) error
	SetRemoteDescription(desc *SessionDescription) error
	AddICECandidate(candidate *ICECandidate) error
	Close() error

	// Connection state
	ConnectionState() ConnectionState
	ICEConnectionState() ICEConnectionState
	SignalingState() SignalingState
	LocalDescription() *SessionDescription
	RemoteDescription() *SessionDescription

	// Media management
	AddTrack(track *MediaStreamTrack) error
	RemoveTrack(track *MediaStreamTrack) error
	GetTracks() []*MediaStreamTrack
	GetLocalTracks() []*MediaStreamTrack
	GetRemoteTracks() []*MediaStreamTrack

	// Data channels
	CreateDataChannel(label string, config *DataChannelConfig) (DataChannel, error)

	// Event handlers
	OnConnectionStateChange(handler func(ConnectionState))
	OnICEConnectionStateChange(handler func(ICEConnectionState))
	OnSignalingStateChange(handler func(SignalingState))
	OnICECandidate(handler func(*ICECandidate))
	OnTrack(handler func(*MediaStreamTrack))
	OnDataChannel(handler func(DataChannel))
	OnError(handler func(error))

	// Statistics
	GetStats() (*PeerConnectionStats, error)

	// Configuration
	GetConfiguration() *PeerConnectionConfig
	SetConfiguration(config *PeerConnectionConfig) error

	// Peer info
	ID() string
	RemotePeerID() string
	SetRemotePeerID(id string)
}

// DataChannel interface định nghĩa data channel
type DataChannel interface {
	// Channel info
	Label() string
	ID() uint16
	Protocol() string
	State() DataChannelState

	// Data transfer
	Send(data []byte) error
	SendText(text string) error
	SendJSON(v interface{}) error

	// Stream interface
	io.ReadWriteCloser

	// Event handlers
	OnOpen(handler func())
	OnClose(handler func())
	OnMessage(handler func([]byte))
	OnError(handler func(error))

	// Configuration
	Ordered() bool
	MaxPacketLifeTime() uint16
	MaxRetransmits() uint16

	// Statistics
	BufferedAmount() uint64
	BufferedAmountLowThreshold() uint64
	SetBufferedAmountLowThreshold(threshold uint64)

	// Lifecycle
	Close() error
}

// MediaEngine interface định nghĩa media engine
type MediaEngine interface {
	// Codec management
	RegisterCodec(codec *Codec) error
	GetCodecs() []*Codec
	GetCodecByName(name string) (*Codec, error)

	// Media processing
	CreateTrack(kind MediaType, id, label string) (*MediaStreamTrack, error)
	CreateLocalTrack(kind MediaType, source MediaSource) (*MediaStreamTrack, error)

	// Stream management
	CreateMediaStream(label string) (*MediaStream, error)
	AddTrackToStream(stream *MediaStream, track *MediaStreamTrack) error
	RemoveTrackFromStream(stream *MediaStream, track *MediaStreamTrack) error

	// Media capture
	GetUserMedia(constraints *MediaConstraints) (*MediaStream, error)
	GetDisplayMedia(constraints *DisplayMediaConstraints) (*MediaStream, error)

	// Media recording
	CreateRecorder(stream *MediaStream, options *RecorderOptions) (MediaRecorder, error)
}

// MediaRecorder interface định nghĩa media recorder
type MediaRecorder interface {
	// Recording control
	Start() error
	Stop() error
	Pause() error
	Resume() error

	// State
	State() RecorderState
	MimeType() string

	// Event handlers
	OnDataAvailable(handler func([]byte))
	OnStart(handler func())
	OnStop(handler func())
	OnPause(handler func())
	OnResume(handler func())
	OnError(handler func(error))

	// Configuration
	SetBitrate(bitrate uint32) error
	SetFramerate(framerate uint32) error

	// Lifecycle
	Close() error
}

// SignalingClient interface định nghĩa signaling client
type SignalingClient interface {
	// Connection management
	Connect(url string) error
	ConnectWithContext(ctx context.Context, url string) error
	Disconnect() error
	IsConnected() bool

	// Message handling
	SendMessage(msg *SignalingMessage) error
	SendOffer(to string, offer *SessionDescription) error
	SendAnswer(to string, answer *SessionDescription) error
	SendICECandidate(to string, candidate *ICECandidate) error
	SendBye(to string) error

	// Room management
	JoinRoom(roomID string, userInfo *PeerInfo) error
	LeaveRoom(roomID string) error
	GetRoomInfo(roomID string) (*RoomInfo, error)
	ListRooms() ([]*RoomInfo, error)

	// Event handlers
	OnMessage(handler func(*SignalingMessage))
	OnOffer(handler func(from string, offer *SessionDescription))
	OnAnswer(handler func(from string, answer *SessionDescription))
	OnICECandidate(handler func(from string, candidate *ICECandidate))
	OnPeerJoined(handler func(*PeerInfo))
	OnPeerLeft(handler func(string))
	OnRoomUpdate(handler func(*RoomInfo))
	OnError(handler func(error))

	// Authentication
	SetAuthToken(token string)
	GetAuthToken() string

	// Configuration
	SetReconnectOptions(enabled bool, interval, maxAttempts int)

	// Lifecycle
	Close() error
}

// SignalingServer interface định nghĩa signaling server
type SignalingServer interface {
	// Server lifecycle
	Start(addr string) error
	StartWithContext(ctx context.Context, addr string) error
	Stop() error
	Shutdown(ctx context.Context) error

	// Room management
	CreateRoom(info *RoomInfo) error
	DeleteRoom(roomID string) error
	GetRoom(roomID string) (*RoomInfo, error)
	ListRooms() ([]*RoomInfo, error)
	GetRoomPeers(roomID string) ([]*PeerInfo, error)

	// Peer management
	GetPeer(peerID string) (*PeerInfo, error)
	DisconnectPeer(peerID string) error
	BroadcastToRoom(roomID string, msg *SignalingMessage) error
	SendToPeer(peerID string, msg *SignalingMessage) error

	// Event handlers
	OnPeerConnected(handler func(*PeerInfo))
	OnPeerDisconnected(handler func(string))
	OnRoomCreated(handler func(*RoomInfo))
	OnRoomDeleted(handler func(string))
	OnMessage(handler func(*SignalingMessage))
	OnError(handler func(error))

	// Authentication
	SetAuthHandler(handler func(token string) (*PeerInfo, error))

	// Middleware
	AddMiddleware(middleware SignalingMiddleware)

	// Statistics
	GetStats() *ServerStats

	// Configuration
	SetConfig(config *ServerConfig)
	GetConfig() *ServerConfig
}

// Room interface định nghĩa room management
type Room interface {
	// Room info
	ID() string
	Info() *RoomInfo

	// Peer management
	AddPeer(peer *PeerInfo) error
	RemovePeer(peerID string) error
	GetPeer(peerID string) (*PeerInfo, error)
	GetPeers() []*PeerInfo
	GetPeerCount() int
	HasPeer(peerID string) bool

	// Broadcasting
	Broadcast(msg *SignalingMessage) error
	BroadcastExcept(msg *SignalingMessage, excludePeerID string) error
	SendToPeer(peerID string, msg *SignalingMessage) error

	// Room state
	IsEmpty() bool
	IsFull() bool
	IsActive() bool

	// Event handlers
	OnPeerJoined(handler func(*PeerInfo))
	OnPeerLeft(handler func(string))
	OnMessage(handler func(*SignalingMessage))
	OnEmpty(handler func())

	// Configuration
	UpdateInfo(info *RoomInfo) error
	SetMaxPeers(max int) error

	// Lifecycle
	Close() error
}

// SFU interface định nghĩa Selective Forwarding Unit
type SFU interface {
	// SFU lifecycle
	Start() error
	Stop() error

	// Room management
	CreateRoom(roomID string, config *RoomConfig) (SFURoom, error)
	GetRoom(roomID string) (SFURoom, error)
	DeleteRoom(roomID string) error

	// Peer management
	AddPeer(roomID, peerID string, pc PeerConnection) error
	RemovePeer(roomID, peerID string) error

	// Media routing
	RouteMedia(fromPeer, toPeer string, track *MediaStreamTrack) error
	StopRouting(fromPeer, toPeer string, trackID string) error

	// Quality control
	SetBitrate(peerID string, bitrate uint32) error
	SetFramerate(peerID string, framerate uint32) error
	SetResolution(peerID string, width, height uint32) error

	// Statistics
	GetRoomStats(roomID string) (*RoomStats, error)
	GetPeerStats(roomID, peerID string) (*PeerStats, error)

	// Event handlers
	OnRoomCreated(handler func(string))
	OnRoomDeleted(handler func(string))
	OnPeerJoined(handler func(string, string))
	OnPeerLeft(handler func(string, string))
	OnTrackAdded(handler func(string, string, *MediaStreamTrack))
	OnTrackRemoved(handler func(string, string, string))

	// Configuration
	SetConfig(config *SFUConfig)
	GetConfig() *SFUConfig
}

// SFURoom interface định nghĩa SFU room
type SFURoom interface {
	// Room info
	ID() string
	Config() *RoomConfig

	// Peer management
	AddPeer(peerID string, pc PeerConnection) error
	RemovePeer(peerID string) error
	GetPeers() []string
	HasPeer(peerID string) bool

	// Media forwarding
	ForwardTrack(fromPeer string, track *MediaStreamTrack) error
	StopForwarding(fromPeer string, trackID string) error

	// Quality adaptation
	AdaptQuality(peerID string, constraints *QualityConstraints) error

	// Statistics
	GetStats() (*RoomStats, error)
	GetPeerStats(peerID string) (*PeerStats, error)

	// Lifecycle
	Close() error
}

// Additional types for interfaces
type OfferOptions struct {
	VoiceActivityDetection bool `json:"voiceActivityDetection"`
	ICERestart             bool `json:"iceRestart"`
}

type AnswerOptions struct {
	VoiceActivityDetection bool `json:"voiceActivityDetection"`
}

type Codec struct {
	Name        string            `json:"name"`
	ClockRate   uint32            `json:"clockRate"`
	Channels    uint16            `json:"channels,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	PayloadType uint8             `json:"payloadType"`
}

type MediaSource interface {
	Read() ([]byte, error)
	Close() error
}

type MediaConstraints struct {
	Audio *AudioConstraints `json:"audio,omitempty"`
	Video *VideoConstraints `json:"video,omitempty"`
}

type AudioConstraints struct {
	Enabled          bool   `json:"enabled"`
	DeviceID         string `json:"deviceId,omitempty"`
	SampleRate       uint32 `json:"sampleRate,omitempty"`
	ChannelCount     uint16 `json:"channelCount,omitempty"`
	EchoCancellation bool   `json:"echoCancellation"`
	NoiseSuppression bool   `json:"noiseSuppression"`
	AutoGainControl  bool   `json:"autoGainControl"`
}

type VideoConstraints struct {
	Enabled   bool   `json:"enabled"`
	DeviceID  string `json:"deviceId,omitempty"`
	Width     uint32 `json:"width,omitempty"`
	Height    uint32 `json:"height,omitempty"`
	Framerate uint32 `json:"framerate,omitempty"`
	Bitrate   uint32 `json:"bitrate,omitempty"`
}

type DisplayMediaConstraints struct {
	Video *VideoConstraints `json:"video,omitempty"`
	Audio *AudioConstraints `json:"audio,omitempty"`
}

type RecorderOptions struct {
	MimeType string `json:"mimeType,omitempty"`
	Bitrate  uint32 `json:"bitrate,omitempty"`
}

type RecorderState int

const (
	RecorderStateInactive RecorderState = iota
	RecorderStateRecording
	RecorderStatePaused
)

type SignalingMiddleware func(next SignalingHandler) SignalingHandler
type SignalingHandler func(*SignalingMessage) error

type ServerStats struct {
	ActiveConnections int64     `json:"activeConnections"`
	TotalRooms        int       `json:"totalRooms"`
	TotalPeers        int       `json:"totalPeers"`
	MessagesPerSecond float64   `json:"messagesPerSecond"`
	Uptime            int64     `json:"uptime"`
	LastUpdated       time.Time `json:"lastUpdated"`
}

type ServerConfig struct {
	MaxRooms        int           `json:"maxRooms"`
	MaxPeersPerRoom int           `json:"maxPeersPerRoom"`
	MessageTimeout  time.Duration `json:"messageTimeout"`
	PeerTimeout     time.Duration `json:"peerTimeout"`
	EnableAuth      bool          `json:"enableAuth"`
	EnableCORS      bool          `json:"enableCORS"`
	AllowedOrigins  []string      `json:"allowedOrigins"`
}

type RoomConfig struct {
	MaxPeers     int    `json:"maxPeers"`
	AudioCodec   string `json:"audioCodec,omitempty"`
	VideoCodec   string `json:"videoCodec,omitempty"`
	MaxBitrate   uint32 `json:"maxBitrate,omitempty"`
	MaxFramerate uint32 `json:"maxFramerate,omitempty"`
}

type SFUConfig struct {
	MaxRooms        int           `json:"maxRooms"`
	MaxPeersPerRoom int           `json:"maxPeersPerRoom"`
	BitrateLimit    uint32        `json:"bitrateLimit"`
	StatsInterval   time.Duration `json:"statsInterval"`
}

type QualityConstraints struct {
	MaxBitrate   uint32 `json:"maxBitrate,omitempty"`
	MaxFramerate uint32 `json:"maxFramerate,omitempty"`
	MaxWidth     uint32 `json:"maxWidth,omitempty"`
	MaxHeight    uint32 `json:"maxHeight,omitempty"`
}

type RoomStats struct {
	RoomID      string    `json:"roomId"`
	PeerCount   int       `json:"peerCount"`
	TotalTracks int       `json:"totalTracks"`
	Bitrate     uint32    `json:"bitrate"`
	PacketLoss  float64   `json:"packetLoss"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type PeerStats struct {
	PeerID      string        `json:"peerId"`
	TrackCount  int           `json:"trackCount"`
	Bitrate     uint32        `json:"bitrate"`
	PacketLoss  float64       `json:"packetLoss"`
	RTT         time.Duration `json:"rtt"`
	LastUpdated time.Time     `json:"lastUpdated"`
}
