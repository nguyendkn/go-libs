package webrtc

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// mediaRecorder implements the MediaRecorder interface
type mediaRecorder struct {
	stream  *MediaStream
	options *RecorderOptions
	
	// State
	state int32 // atomic RecorderState
	
	// Configuration
	mimeType  string
	bitrate   uint32
	framerate uint32
	
	// Event handlers
	onDataAvailable func([]byte)
	onStart         func()
	onStop          func()
	onPause         func()
	onResume        func()
	onError         func(error)
	handlersMu      sync.RWMutex
	
	// Recording control
	recording chan struct{}
	paused    chan struct{}
	stopped   chan struct{}
	
	// Data buffer
	buffer []byte
	bufMu  sync.Mutex
}

// newMediaRecorder tạo một MediaRecorder mới
func newMediaRecorder(stream *MediaStream, options *RecorderOptions) (MediaRecorder, error) {
	if stream == nil {
		return nil, fmt.Errorf("stream cannot be nil")
	}
	
	if options == nil {
		options = &RecorderOptions{
			MimeType: "video/webm",
			Bitrate:  1000000, // 1 Mbps
		}
	}
	
	recorder := &mediaRecorder{
		stream:    stream,
		options:   options,
		mimeType:  options.MimeType,
		bitrate:   options.Bitrate,
		framerate: 30, // default framerate
		recording: make(chan struct{}),
		paused:    make(chan struct{}),
		stopped:   make(chan struct{}),
	}
	
	atomic.StoreInt32(&recorder.state, int32(RecorderStateInactive))
	
	return recorder, nil
}

// Recording control
func (mr *mediaRecorder) Start() error {
	if !atomic.CompareAndSwapInt32(&mr.state, int32(RecorderStateInactive), int32(RecorderStateRecording)) {
		return fmt.Errorf("recorder is not in inactive state")
	}
	
	// Start recording goroutine
	go mr.recordingLoop()
	
	mr.handlersMu.RLock()
	if mr.onStart != nil {
		go mr.onStart()
	}
	mr.handlersMu.RUnlock()
	
	return nil
}

func (mr *mediaRecorder) Stop() error {
	currentState := atomic.LoadInt32(&mr.state)
	if currentState == int32(RecorderStateInactive) {
		return fmt.Errorf("recorder is not recording")
	}
	
	atomic.StoreInt32(&mr.state, int32(RecorderStateInactive))
	close(mr.stopped)
	
	// Flush any remaining data
	mr.flushBuffer()
	
	mr.handlersMu.RLock()
	if mr.onStop != nil {
		go mr.onStop()
	}
	mr.handlersMu.RUnlock()
	
	return nil
}

func (mr *mediaRecorder) Pause() error {
	if !atomic.CompareAndSwapInt32(&mr.state, int32(RecorderStateRecording), int32(RecorderStatePaused)) {
		return fmt.Errorf("recorder is not recording")
	}
	
	mr.handlersMu.RLock()
	if mr.onPause != nil {
		go mr.onPause()
	}
	mr.handlersMu.RUnlock()
	
	return nil
}

func (mr *mediaRecorder) Resume() error {
	if !atomic.CompareAndSwapInt32(&mr.state, int32(RecorderStatePaused), int32(RecorderStateRecording)) {
		return fmt.Errorf("recorder is not paused")
	}
	
	mr.handlersMu.RLock()
	if mr.onResume != nil {
		go mr.onResume()
	}
	mr.handlersMu.RUnlock()
	
	return nil
}

// State
func (mr *mediaRecorder) State() RecorderState {
	return RecorderState(atomic.LoadInt32(&mr.state))
}

func (mr *mediaRecorder) MimeType() string {
	return mr.mimeType
}

// Event handlers
func (mr *mediaRecorder) OnDataAvailable(handler func([]byte)) {
	mr.handlersMu.Lock()
	mr.onDataAvailable = handler
	mr.handlersMu.Unlock()
}

func (mr *mediaRecorder) OnStart(handler func()) {
	mr.handlersMu.Lock()
	mr.onStart = handler
	mr.handlersMu.Unlock()
}

func (mr *mediaRecorder) OnStop(handler func()) {
	mr.handlersMu.Lock()
	mr.onStop = handler
	mr.handlersMu.Unlock()
}

func (mr *mediaRecorder) OnPause(handler func()) {
	mr.handlersMu.Lock()
	mr.onPause = handler
	mr.handlersMu.Unlock()
}

func (mr *mediaRecorder) OnResume(handler func()) {
	mr.handlersMu.Lock()
	mr.onResume = handler
	mr.handlersMu.Unlock()
}

func (mr *mediaRecorder) OnError(handler func(error)) {
	mr.handlersMu.Lock()
	mr.onError = handler
	mr.handlersMu.Unlock()
}

// Configuration
func (mr *mediaRecorder) SetBitrate(bitrate uint32) error {
	mr.bitrate = bitrate
	return nil
}

func (mr *mediaRecorder) SetFramerate(framerate uint32) error {
	mr.framerate = framerate
	return nil
}

// Lifecycle
func (mr *mediaRecorder) Close() error {
	if mr.State() != RecorderStateInactive {
		mr.Stop()
	}
	return nil
}

// recordingLoop chạy vòng lặp recording chính
func (mr *mediaRecorder) recordingLoop() {
	ticker := time.NewTicker(time.Second / time.Duration(mr.framerate))
	defer ticker.Stop()
	
	for {
		select {
		case <-mr.stopped:
			return
			
		case <-ticker.C:
			if mr.State() == RecorderStateRecording {
				mr.captureFrame()
			}
		}
	}
}

// captureFrame capture một frame từ stream
func (mr *mediaRecorder) captureFrame() {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Read video/audio data from the stream tracks
	// 2. Encode the data according to the mime type
	// 3. Add to buffer or emit data available event
	
	// Simulate capturing data
	data := make([]byte, 1024) // Dummy data
	for i := range data {
		data[i] = byte(time.Now().UnixNano() % 256)
	}
	
	mr.bufMu.Lock()
	mr.buffer = append(mr.buffer, data...)
	
	// Emit data if buffer is large enough
	if len(mr.buffer) >= 8192 { // 8KB chunks
		mr.emitData()
	}
	mr.bufMu.Unlock()
}

// emitData emit buffered data
func (mr *mediaRecorder) emitData() {
	if len(mr.buffer) == 0 {
		return
	}
	
	data := make([]byte, len(mr.buffer))
	copy(data, mr.buffer)
	mr.buffer = mr.buffer[:0] // Clear buffer
	
	mr.handlersMu.RLock()
	if mr.onDataAvailable != nil {
		go mr.onDataAvailable(data)
	}
	mr.handlersMu.RUnlock()
}

// flushBuffer flush remaining data in buffer
func (mr *mediaRecorder) flushBuffer() {
	mr.bufMu.Lock()
	defer mr.bufMu.Unlock()
	
	if len(mr.buffer) > 0 {
		mr.emitData()
	}
}
