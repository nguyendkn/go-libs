package webrtc

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pion/webrtc/v4"
)

// dataChannel implements the DataChannel interface
type dataChannel struct {
	dc *webrtc.DataChannel
	
	// State
	state int32 // atomic DataChannelState
	
	// Event handlers
	onOpen    func()
	onClose   func()
	onMessage func([]byte)
	onError   func(error)
	mu        sync.RWMutex
	
	// Buffer management
	bufferedAmountLowThreshold uint64
}

// newDataChannel tạo một DataChannel mới từ Pion DataChannel
func newDataChannel(dc *webrtc.DataChannel) DataChannel {
	channel := &dataChannel{
		dc: dc,
	}
	
	// Set initial state
	atomic.StoreInt32(&channel.state, int32(DataChannelStateConnecting))
	
	// Setup Pion event handlers
	channel.setupPionHandlers()
	
	return channel
}

// setupPionHandlers thiết lập event handlers cho Pion DataChannel
func (dc *dataChannel) setupPionHandlers() {
	// OnOpen
	dc.dc.OnOpen(func() {
		atomic.StoreInt32(&dc.state, int32(DataChannelStateOpen))
		
		dc.mu.RLock()
		if dc.onOpen != nil {
			go dc.onOpen()
		}
		dc.mu.RUnlock()
	})
	
	// OnClose
	dc.dc.OnClose(func() {
		atomic.StoreInt32(&dc.state, int32(DataChannelStateClosed))
		
		dc.mu.RLock()
		if dc.onClose != nil {
			go dc.onClose()
		}
		dc.mu.RUnlock()
	})
	
	// OnMessage
	dc.dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		dc.mu.RLock()
		if dc.onMessage != nil {
			go dc.onMessage(msg.Data)
		}
		dc.mu.RUnlock()
	})
	
	// OnError
	dc.dc.OnError(func(err error) {
		dc.mu.RLock()
		if dc.onError != nil {
			go dc.onError(err)
		}
		dc.mu.RUnlock()
	})
}

// Channel info methods
func (dc *dataChannel) Label() string {
	return dc.dc.Label()
}

func (dc *dataChannel) ID() uint16 {
	if dc.dc.ID() != nil {
		return *dc.dc.ID()
	}
	return 0
}

func (dc *dataChannel) Protocol() string {
	return dc.dc.Protocol()
}

func (dc *dataChannel) State() DataChannelState {
	return DataChannelState(atomic.LoadInt32(&dc.state))
}

// Data transfer methods
func (dc *dataChannel) Send(data []byte) error {
	if dc.State() != DataChannelStateOpen {
		return ErrDataChannelClosed
	}
	
	return dc.dc.Send(data)
}

func (dc *dataChannel) SendText(text string) error {
	return dc.Send([]byte(text))
}

func (dc *dataChannel) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return dc.Send(data)
}

// io.ReadWriteCloser implementation
func (dc *dataChannel) Read(p []byte) (n int, err error) {
	// This is a simplified implementation
	// In a real implementation, you would need to buffer incoming messages
	return 0, fmt.Errorf("read not implemented for data channel")
}

func (dc *dataChannel) Write(p []byte) (n int, err error) {
	err = dc.Send(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (dc *dataChannel) Close() error {
	atomic.StoreInt32(&dc.state, int32(DataChannelStateClosing))
	return dc.dc.Close()
}

// Event handler setters
func (dc *dataChannel) OnOpen(handler func()) {
	dc.mu.Lock()
	dc.onOpen = handler
	dc.mu.Unlock()
}

func (dc *dataChannel) OnClose(handler func()) {
	dc.mu.Lock()
	dc.onClose = handler
	dc.mu.Unlock()
}

func (dc *dataChannel) OnMessage(handler func([]byte)) {
	dc.mu.Lock()
	dc.onMessage = handler
	dc.mu.Unlock()
}

func (dc *dataChannel) OnError(handler func(error)) {
	dc.mu.Lock()
	dc.onError = handler
	dc.mu.Unlock()
}

// Configuration methods
func (dc *dataChannel) Ordered() bool {
	return dc.dc.Ordered()
}

func (dc *dataChannel) MaxPacketLifeTime() uint16 {
	if dc.dc.MaxPacketLifeTime() != nil {
		return *dc.dc.MaxPacketLifeTime()
	}
	return 0
}

func (dc *dataChannel) MaxRetransmits() uint16 {
	if dc.dc.MaxRetransmits() != nil {
		return *dc.dc.MaxRetransmits()
	}
	return 0
}

// Buffer management
func (dc *dataChannel) BufferedAmount() uint64 {
	return dc.dc.BufferedAmount()
}

func (dc *dataChannel) BufferedAmountLowThreshold() uint64 {
	return dc.bufferedAmountLowThreshold
}

func (dc *dataChannel) SetBufferedAmountLowThreshold(threshold uint64) {
	dc.bufferedAmountLowThreshold = threshold
	dc.dc.SetBufferedAmountLowThreshold(threshold)
}
