package webrtc

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
)

// peerConnection implements the PeerConnection interface
type peerConnection struct {
	id           string
	remotePeerID string

	// Pion WebRTC peer connection
	pc *webrtc.PeerConnection

	// Configuration
	config *PeerConnectionConfig

	// State management
	connectionState    int32 // atomic ConnectionState
	iceConnectionState int32 // atomic ICEConnectionState
	signalingState     int32 // atomic SignalingState

	// Tracks
	localTracks  map[string]*MediaStreamTrack
	remoteTracks map[string]*MediaStreamTrack
	tracksMu     sync.RWMutex

	// Data channels
	dataChannels map[string]DataChannel
	channelsMu   sync.RWMutex

	// Event handlers
	onConnectionStateChange    func(ConnectionState)
	onICEConnectionStateChange func(ICEConnectionState)
	onSignalingStateChange     func(SignalingState)
	onICECandidate             func(*ICECandidate)
	onTrack                    func(*MediaStreamTrack)
	onDataChannel              func(DataChannel)
	onError                    func(error)
	handlersMu                 sync.RWMutex

	// Statistics
	stats     *PeerConnectionStats
	statsMu   sync.RWMutex
	statsStop chan struct{}

	// Context and lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	closed int32 // atomic

	// Wait group for goroutines
	wg sync.WaitGroup
}

// NewPeerConnection tạo một PeerConnection mới
func NewPeerConnection(config *PeerConnectionConfig) (PeerConnection, error) {
	if config == nil {
		config = &PeerConnectionConfig{
			ICEServers:          DefaultICEServers,
			ConnectionTimeout:   DefaultConnectionTimeout,
			DisconnectedTimeout: DefaultDisconnectedTimeout,
			FailedTimeout:       DefaultFailedTimeout,
			KeepAliveInterval:   DefaultKeepAliveInterval,
		}
	}

	// Convert to Pion WebRTC config
	pionConfig := webrtc.Configuration{
		ICEServers: make([]webrtc.ICEServer, len(config.ICEServers)),
	}

	for i, server := range config.ICEServers {
		pionConfig.ICEServers[i] = webrtc.ICEServer{
			URLs:       server.URLs,
			Username:   server.Username,
			Credential: server.Credential,
		}
	}

	// Create Pion peer connection
	pc, err := webrtc.NewPeerConnection(pionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn := &peerConnection{
		id:           uuid.New().String(),
		pc:           pc,
		config:       config,
		localTracks:  make(map[string]*MediaStreamTrack),
		remoteTracks: make(map[string]*MediaStreamTrack),
		dataChannels: make(map[string]DataChannel),
		stats: &PeerConnectionStats{
			ConnectionState:    ConnectionStateNew,
			ICEConnectionState: ICEConnectionStateNew,
			SignalingState:     SignalingStateStable,
			ConnectedAt:        time.Now(),
			LastActivity:       time.Now(),
		},
		statsStop: make(chan struct{}),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Set initial states
	atomic.StoreInt32(&conn.connectionState, int32(ConnectionStateNew))
	atomic.StoreInt32(&conn.iceConnectionState, int32(ICEConnectionStateNew))
	atomic.StoreInt32(&conn.signalingState, int32(SignalingStateStable))

	// Setup Pion event handlers
	conn.setupPionHandlers()

	// Start statistics collection
	conn.wg.Add(1)
	go conn.collectStats()

	return conn, nil
}

// setupPionHandlers thiết lập event handlers cho Pion WebRTC
func (pc *peerConnection) setupPionHandlers() {
	// Connection state change
	pc.pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		var newState ConnectionState
		switch state {
		case webrtc.PeerConnectionStateNew:
			newState = ConnectionStateNew
		case webrtc.PeerConnectionStateConnecting:
			newState = ConnectionStateConnecting
		case webrtc.PeerConnectionStateConnected:
			newState = ConnectionStateConnected
		case webrtc.PeerConnectionStateDisconnected:
			newState = ConnectionStateDisconnected
		case webrtc.PeerConnectionStateFailed:
			newState = ConnectionStateFailed
		case webrtc.PeerConnectionStateClosed:
			newState = ConnectionStateClosed
		}

		atomic.StoreInt32(&pc.connectionState, int32(newState))

		pc.statsMu.Lock()
		pc.stats.ConnectionState = newState
		pc.stats.LastActivity = time.Now()
		if newState == ConnectionStateConnected {
			pc.stats.ConnectedAt = time.Now()
		}
		pc.statsMu.Unlock()

		pc.handlersMu.RLock()
		if pc.onConnectionStateChange != nil {
			go pc.onConnectionStateChange(newState)
		}
		pc.handlersMu.RUnlock()
	})

	// ICE connection state change
	pc.pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		var newState ICEConnectionState
		switch state {
		case webrtc.ICEConnectionStateNew:
			newState = ICEConnectionStateNew
		case webrtc.ICEConnectionStateChecking:
			newState = ICEConnectionStateChecking
		case webrtc.ICEConnectionStateConnected:
			newState = ICEConnectionStateConnected
		case webrtc.ICEConnectionStateCompleted:
			newState = ICEConnectionStateCompleted
		case webrtc.ICEConnectionStateDisconnected:
			newState = ICEConnectionStateDisconnected
		case webrtc.ICEConnectionStateFailed:
			newState = ICEConnectionStateFailed
		case webrtc.ICEConnectionStateClosed:
			newState = ICEConnectionStateClosed
		}

		atomic.StoreInt32(&pc.iceConnectionState, int32(newState))

		pc.statsMu.Lock()
		pc.stats.ICEConnectionState = newState
		pc.stats.LastActivity = time.Now()
		pc.statsMu.Unlock()

		pc.handlersMu.RLock()
		if pc.onICEConnectionStateChange != nil {
			go pc.onICEConnectionStateChange(newState)
		}
		pc.handlersMu.RUnlock()
	})

	// Signaling state change
	pc.pc.OnSignalingStateChange(func(state webrtc.SignalingState) {
		var newState SignalingState
		switch state {
		case webrtc.SignalingStateStable:
			newState = SignalingStateStable
		case webrtc.SignalingStateHaveLocalOffer:
			newState = SignalingStateHaveLocalOffer
		case webrtc.SignalingStateHaveRemoteOffer:
			newState = SignalingStateHaveRemoteOffer
		case webrtc.SignalingStateHaveLocalPranswer:
			newState = SignalingStateHaveLocalPranswer
		case webrtc.SignalingStateHaveRemotePranswer:
			newState = SignalingStateHaveRemotePranswer
		case webrtc.SignalingStateClosed:
			newState = SignalingStateClosed
		}

		atomic.StoreInt32(&pc.signalingState, int32(newState))

		pc.statsMu.Lock()
		pc.stats.SignalingState = newState
		pc.stats.LastActivity = time.Now()
		pc.statsMu.Unlock()

		pc.handlersMu.RLock()
		if pc.onSignalingStateChange != nil {
			go pc.onSignalingStateChange(newState)
		}
		pc.handlersMu.RUnlock()
	})

	// ICE candidate
	pc.pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		iceCandidate := &ICECandidate{
			Candidate:     candidate.String(),
			SDPMid:        "",
			SDPMLineIndex: 0,
		}

		pc.handlersMu.RLock()
		if pc.onICECandidate != nil {
			go pc.onICECandidate(iceCandidate)
		}
		pc.handlersMu.RUnlock()
	})

	// Track received
	pc.pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		mediaTrack := &MediaStreamTrack{
			ID:         track.ID(),
			Label:      track.ID(), // Use ID as label since Label() doesn't exist
			Enabled:    true,
			Muted:      false,
			ReadyState: "live",
			Direction:  TrackDirectionRecvOnly,
			TrackRef:   track,
		}

		// Determine media type
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			mediaTrack.Kind = MediaTypeAudio
		case webrtc.RTPCodecTypeVideo:
			mediaTrack.Kind = MediaTypeVideo
		}

		pc.tracksMu.Lock()
		pc.remoteTracks[track.ID()] = mediaTrack
		pc.tracksMu.Unlock()

		pc.handlersMu.RLock()
		if pc.onTrack != nil {
			go pc.onTrack(mediaTrack)
		}
		pc.handlersMu.RUnlock()
	})

	// Data channel received
	pc.pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dataChannel := newDataChannel(dc)

		pc.channelsMu.Lock()
		pc.dataChannels[dc.Label()] = dataChannel
		pc.channelsMu.Unlock()

		pc.handlersMu.RLock()
		if pc.onDataChannel != nil {
			go pc.onDataChannel(dataChannel)
		}
		pc.handlersMu.RUnlock()
	})
}

// CreateOffer tạo offer SDP
func (pc *peerConnection) CreateOffer(options *OfferOptions) (*SessionDescription, error) {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return nil, ErrPeerConnectionClosed
	}

	var pionOptions *webrtc.OfferOptions
	if options != nil {
		pionOptions = &webrtc.OfferOptions{
			ICERestart: options.ICERestart,
		}
	}

	offer, err := pc.pc.CreateOffer(pionOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	return &SessionDescription{
		Type: offer.Type.String(),
		SDP:  offer.SDP,
	}, nil
}

// CreateAnswer tạo answer SDP
func (pc *peerConnection) CreateAnswer(options *AnswerOptions) (*SessionDescription, error) {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return nil, ErrPeerConnectionClosed
	}

	var pionOptions *webrtc.AnswerOptions
	if options != nil {
		pionOptions = &webrtc.AnswerOptions{}
	}

	answer, err := pc.pc.CreateAnswer(pionOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	return &SessionDescription{
		Type: answer.Type.String(),
		SDP:  answer.SDP,
	}, nil
}

// SetLocalDescription set local SDP
func (pc *peerConnection) SetLocalDescription(desc *SessionDescription) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	sdpType := webrtc.SDPTypeOffer
	switch desc.Type {
	case "offer":
		sdpType = webrtc.SDPTypeOffer
	case "answer":
		sdpType = webrtc.SDPTypeAnswer
	case "pranswer":
		sdpType = webrtc.SDPTypePranswer
	case "rollback":
		sdpType = webrtc.SDPTypeRollback
	default:
		return fmt.Errorf("invalid SDP type: %s", desc.Type)
	}

	sessionDesc := webrtc.SessionDescription{
		Type: sdpType,
		SDP:  desc.SDP,
	}

	if err := pc.pc.SetLocalDescription(sessionDesc); err != nil {
		return fmt.Errorf("failed to set local description: %w", err)
	}

	return nil
}

// SetRemoteDescription set remote SDP
func (pc *peerConnection) SetRemoteDescription(desc *SessionDescription) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	sdpType := webrtc.SDPTypeOffer
	switch desc.Type {
	case "offer":
		sdpType = webrtc.SDPTypeOffer
	case "answer":
		sdpType = webrtc.SDPTypeAnswer
	case "pranswer":
		sdpType = webrtc.SDPTypePranswer
	case "rollback":
		sdpType = webrtc.SDPTypeRollback
	default:
		return fmt.Errorf("invalid SDP type: %s", desc.Type)
	}

	sessionDesc := webrtc.SessionDescription{
		Type: sdpType,
		SDP:  desc.SDP,
	}

	if err := pc.pc.SetRemoteDescription(sessionDesc); err != nil {
		return fmt.Errorf("failed to set remote description: %w", err)
	}

	return nil
}

// AddICECandidate thêm ICE candidate
func (pc *peerConnection) AddICECandidate(candidate *ICECandidate) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	iceCandidate := webrtc.ICECandidateInit{
		Candidate:     candidate.Candidate,
		SDPMid:        &candidate.SDPMid,
		SDPMLineIndex: &candidate.SDPMLineIndex,
	}

	if err := pc.pc.AddICECandidate(iceCandidate); err != nil {
		return fmt.Errorf("failed to add ICE candidate: %w", err)
	}

	return nil
}

// Close đóng peer connection
func (pc *peerConnection) Close() error {
	if !atomic.CompareAndSwapInt32(&pc.closed, 0, 1) {
		return nil // already closed
	}

	pc.cancel()

	// Close statistics collection
	close(pc.statsStop)

	// Close all data channels
	pc.channelsMu.Lock()
	for _, dc := range pc.dataChannels {
		dc.Close()
	}
	pc.channelsMu.Unlock()

	// Close Pion peer connection
	if err := pc.pc.Close(); err != nil {
		return fmt.Errorf("failed to close peer connection: %w", err)
	}

	// Wait for goroutines to finish
	pc.wg.Wait()

	return nil
}

// Connection state methods
func (pc *peerConnection) ConnectionState() ConnectionState {
	return ConnectionState(atomic.LoadInt32(&pc.connectionState))
}

func (pc *peerConnection) ICEConnectionState() ICEConnectionState {
	return ICEConnectionState(atomic.LoadInt32(&pc.iceConnectionState))
}

func (pc *peerConnection) SignalingState() SignalingState {
	return SignalingState(atomic.LoadInt32(&pc.signalingState))
}

func (pc *peerConnection) LocalDescription() *SessionDescription {
	desc := pc.pc.LocalDescription()
	if desc == nil {
		return nil
	}
	return &SessionDescription{
		Type: desc.Type.String(),
		SDP:  desc.SDP,
	}
}

func (pc *peerConnection) RemoteDescription() *SessionDescription {
	desc := pc.pc.RemoteDescription()
	if desc == nil {
		return nil
	}
	return &SessionDescription{
		Type: desc.Type.String(),
		SDP:  desc.SDP,
	}
}

// Track management methods
func (pc *peerConnection) AddTrack(track *MediaStreamTrack) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	// This is a simplified implementation
	// In a real implementation, you would create a proper track from the MediaStreamTrack
	pc.tracksMu.Lock()
	pc.localTracks[track.ID] = track
	pc.tracksMu.Unlock()

	return nil
}

func (pc *peerConnection) RemoveTrack(track *MediaStreamTrack) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	pc.tracksMu.Lock()
	delete(pc.localTracks, track.ID)
	pc.tracksMu.Unlock()

	return nil
}

func (pc *peerConnection) GetTracks() []*MediaStreamTrack {
	pc.tracksMu.RLock()
	defer pc.tracksMu.RUnlock()

	tracks := make([]*MediaStreamTrack, 0, len(pc.localTracks)+len(pc.remoteTracks))
	for _, track := range pc.localTracks {
		tracks = append(tracks, track)
	}
	for _, track := range pc.remoteTracks {
		tracks = append(tracks, track)
	}

	return tracks
}

func (pc *peerConnection) GetLocalTracks() []*MediaStreamTrack {
	pc.tracksMu.RLock()
	defer pc.tracksMu.RUnlock()

	tracks := make([]*MediaStreamTrack, 0, len(pc.localTracks))
	for _, track := range pc.localTracks {
		tracks = append(tracks, track)
	}

	return tracks
}

func (pc *peerConnection) GetRemoteTracks() []*MediaStreamTrack {
	pc.tracksMu.RLock()
	defer pc.tracksMu.RUnlock()

	tracks := make([]*MediaStreamTrack, 0, len(pc.remoteTracks))
	for _, track := range pc.remoteTracks {
		tracks = append(tracks, track)
	}

	return tracks
}

// CreateDataChannel tạo data channel
func (pc *peerConnection) CreateDataChannel(label string, config *DataChannelConfig) (DataChannel, error) {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return nil, ErrPeerConnectionClosed
	}

	var pionConfig *webrtc.DataChannelInit
	if config != nil {
		pionConfig = &webrtc.DataChannelInit{
			Ordered: &config.Ordered,
		}

		if config.Protocol != "" {
			pionConfig.Protocol = &config.Protocol
		}
		if config.Negotiated {
			pionConfig.Negotiated = &config.Negotiated
		}
		if config.ID != 0 {
			pionConfig.ID = &config.ID
		}
		// Only set one of MaxPacketLifeTime or MaxRetransmits, not both
		if config.MaxPacketLifeTime > 0 {
			pionConfig.MaxPacketLifeTime = &config.MaxPacketLifeTime
		} else if config.MaxRetransmits > 0 {
			pionConfig.MaxRetransmits = &config.MaxRetransmits
		}
	}

	dc, err := pc.pc.CreateDataChannel(label, pionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create data channel: %w", err)
	}

	dataChannel := newDataChannel(dc)

	pc.channelsMu.Lock()
	pc.dataChannels[label] = dataChannel
	pc.channelsMu.Unlock()

	return dataChannel, nil
}

// Event handler setters
func (pc *peerConnection) OnConnectionStateChange(handler func(ConnectionState)) {
	pc.handlersMu.Lock()
	pc.onConnectionStateChange = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnICEConnectionStateChange(handler func(ICEConnectionState)) {
	pc.handlersMu.Lock()
	pc.onICEConnectionStateChange = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnSignalingStateChange(handler func(SignalingState)) {
	pc.handlersMu.Lock()
	pc.onSignalingStateChange = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnICECandidate(handler func(*ICECandidate)) {
	pc.handlersMu.Lock()
	pc.onICECandidate = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnTrack(handler func(*MediaStreamTrack)) {
	pc.handlersMu.Lock()
	pc.onTrack = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnDataChannel(handler func(DataChannel)) {
	pc.handlersMu.Lock()
	pc.onDataChannel = handler
	pc.handlersMu.Unlock()
}

func (pc *peerConnection) OnError(handler func(error)) {
	pc.handlersMu.Lock()
	pc.onError = handler
	pc.handlersMu.Unlock()
}

// GetStats trả về statistics
func (pc *peerConnection) GetStats() (*PeerConnectionStats, error) {
	pc.statsMu.RLock()
	defer pc.statsMu.RUnlock()

	// Create a copy of stats
	stats := *pc.stats
	return &stats, nil
}

// Configuration methods
func (pc *peerConnection) GetConfiguration() *PeerConnectionConfig {
	return pc.config
}

func (pc *peerConnection) SetConfiguration(config *PeerConnectionConfig) error {
	if atomic.LoadInt32(&pc.closed) == 1 {
		return ErrPeerConnectionClosed
	}

	pc.config = config
	return nil
}

// Peer info methods
func (pc *peerConnection) ID() string {
	return pc.id
}

func (pc *peerConnection) RemotePeerID() string {
	return pc.remotePeerID
}

func (pc *peerConnection) SetRemotePeerID(id string) {
	pc.remotePeerID = id
}

// collectStats thu thập statistics định kỳ
func (pc *peerConnection) collectStats() {
	defer pc.wg.Done()

	ticker := time.NewTicker(DefaultStatsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pc.statsStop:
			return
		case <-pc.ctx.Done():
			return
		case <-ticker.C:
			pc.updateStats()
		}
	}
}

// updateStats cập nhật statistics
func (pc *peerConnection) updateStats() {
	// Simplified stats update - in a real implementation,
	// you would parse the Pion WebRTC stats properly
	pc.statsMu.Lock()
	defer pc.statsMu.Unlock()

	// Update last activity
	pc.stats.LastActivity = time.Now()

	// In a real implementation, you would:
	// 1. Get stats from pc.pc.GetStats()
	// 2. Parse the stats map
	// 3. Update the relevant fields
	// For now, we'll just update the timestamp
}
