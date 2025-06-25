package webrtc

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
)

// mediaEngine implements the MediaEngine interface
type mediaEngine struct {
	engine *webrtc.MediaEngine
	
	// Codec registry
	codecs   map[string]*Codec
	codecsMu sync.RWMutex
	
	// Track registry
	tracks   map[string]*MediaStreamTrack
	tracksMu sync.RWMutex
	
	// Stream registry
	streams   map[string]*MediaStream
	streamsMu sync.RWMutex
}

// NewMediaEngine tạo một MediaEngine mới
func NewMediaEngine() MediaEngine {
	engine := &mediaEngine{
		engine:  &webrtc.MediaEngine{},
		codecs:  make(map[string]*Codec),
		tracks:  make(map[string]*MediaStreamTrack),
		streams: make(map[string]*MediaStream),
	}
	
	// Register default codecs
	engine.registerDefaultCodecs()
	
	return engine
}

// registerDefaultCodecs đăng ký các codec mặc định
func (me *mediaEngine) registerDefaultCodecs() {
	// Audio codecs
	me.RegisterCodec(&Codec{
		Name:        "opus",
		ClockRate:   48000,
		Channels:    2,
		PayloadType: 111,
		Parameters: map[string]string{
			"minptime": "10",
			"useinbandfec": "1",
		},
	})
	
	me.RegisterCodec(&Codec{
		Name:        "PCMU",
		ClockRate:   8000,
		Channels:    1,
		PayloadType: 0,
	})
	
	me.RegisterCodec(&Codec{
		Name:        "PCMA",
		ClockRate:   8000,
		Channels:    1,
		PayloadType: 8,
	})
	
	// Video codecs
	me.RegisterCodec(&Codec{
		Name:        "VP8",
		ClockRate:   90000,
		PayloadType: 96,
	})
	
	me.RegisterCodec(&Codec{
		Name:        "VP9",
		ClockRate:   90000,
		PayloadType: 98,
		Parameters: map[string]string{
			"profile-id": "0",
		},
	})
	
	me.RegisterCodec(&Codec{
		Name:        "H264",
		ClockRate:   90000,
		PayloadType: 102,
		Parameters: map[string]string{
			"level-asymmetry-allowed": "1",
			"packetization-mode":      "1",
			"profile-level-id":        "42001f",
		},
	})
}

// Codec management
func (me *mediaEngine) RegisterCodec(codec *Codec) error {
	me.codecsMu.Lock()
	defer me.codecsMu.Unlock()
	
	me.codecs[codec.Name] = codec
	
	// Register with Pion MediaEngine
	switch codec.Name {
	case "opus":
		if err := me.engine.RegisterCodec(webrtc.RTPCodecParameters{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeOpus,
				ClockRate:    codec.ClockRate,
				Channels:     codec.Channels,
				SDPFmtpLine:  formatParameters(codec.Parameters),
			},
			PayloadType: webrtc.PayloadType(codec.PayloadType),
		}, webrtc.RTPCodecTypeAudio); err != nil {
			return fmt.Errorf("failed to register opus codec: %w", err)
		}
		
	case "VP8":
		if err := me.engine.RegisterCodec(webrtc.RTPCodecParameters{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:    webrtc.MimeTypeVP8,
				ClockRate:   codec.ClockRate,
				SDPFmtpLine: formatParameters(codec.Parameters),
			},
			PayloadType: webrtc.PayloadType(codec.PayloadType),
		}, webrtc.RTPCodecTypeVideo); err != nil {
			return fmt.Errorf("failed to register VP8 codec: %w", err)
		}
		
	case "VP9":
		if err := me.engine.RegisterCodec(webrtc.RTPCodecParameters{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:    webrtc.MimeTypeVP9,
				ClockRate:   codec.ClockRate,
				SDPFmtpLine: formatParameters(codec.Parameters),
			},
			PayloadType: webrtc.PayloadType(codec.PayloadType),
		}, webrtc.RTPCodecTypeVideo); err != nil {
			return fmt.Errorf("failed to register VP9 codec: %w", err)
		}
		
	case "H264":
		if err := me.engine.RegisterCodec(webrtc.RTPCodecParameters{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:    webrtc.MimeTypeH264,
				ClockRate:   codec.ClockRate,
				SDPFmtpLine: formatParameters(codec.Parameters),
			},
			PayloadType: webrtc.PayloadType(codec.PayloadType),
		}, webrtc.RTPCodecTypeVideo); err != nil {
			return fmt.Errorf("failed to register H264 codec: %w", err)
		}
	}
	
	return nil
}

func (me *mediaEngine) GetCodecs() []*Codec {
	me.codecsMu.RLock()
	defer me.codecsMu.RUnlock()
	
	codecs := make([]*Codec, 0, len(me.codecs))
	for _, codec := range me.codecs {
		codecs = append(codecs, codec)
	}
	
	return codecs
}

func (me *mediaEngine) GetCodecByName(name string) (*Codec, error) {
	me.codecsMu.RLock()
	defer me.codecsMu.RUnlock()
	
	codec, exists := me.codecs[name]
	if !exists {
		return nil, fmt.Errorf("codec %s not found", name)
	}
	
	return codec, nil
}

// Media processing
func (me *mediaEngine) CreateTrack(kind MediaType, id, label string) (*MediaStreamTrack, error) {
	if id == "" {
		id = uuid.New().String()
	}
	
	track := &MediaStreamTrack{
		ID:         id,
		Kind:       kind,
		Label:      label,
		Enabled:    true,
		Muted:      false,
		ReadyState: "live",
		Direction:  TrackDirectionSendRecv,
	}
	
	me.tracksMu.Lock()
	me.tracks[id] = track
	me.tracksMu.Unlock()
	
	return track, nil
}

func (me *mediaEngine) CreateLocalTrack(kind MediaType, source MediaSource) (*MediaStreamTrack, error) {
	track, err := me.CreateTrack(kind, "", "")
	if err != nil {
		return nil, err
	}
	
	track.Direction = TrackDirectionSendOnly
	
	// In a real implementation, you would:
	// 1. Create a Pion WebRTC track from the source
	// 2. Set up the media pipeline
	// 3. Start reading from the source and writing to the track
	
	return track, nil
}

// Stream management
func (me *mediaEngine) CreateMediaStream(label string) (*MediaStream, error) {
	id := uuid.New().String()
	
	stream := &MediaStream{
		ID:     id,
		Label:  label,
		Tracks: make([]*MediaStreamTrack, 0),
		Active: true,
	}
	
	me.streamsMu.Lock()
	me.streams[id] = stream
	me.streamsMu.Unlock()
	
	return stream, nil
}

func (me *mediaEngine) AddTrackToStream(stream *MediaStream, track *MediaStreamTrack) error {
	me.streamsMu.Lock()
	defer me.streamsMu.Unlock()
	
	// Check if track already exists in stream
	for _, t := range stream.Tracks {
		if t.ID == track.ID {
			return fmt.Errorf("track %s already exists in stream", track.ID)
		}
	}
	
	stream.Tracks = append(stream.Tracks, track)
	return nil
}

func (me *mediaEngine) RemoveTrackFromStream(stream *MediaStream, track *MediaStreamTrack) error {
	me.streamsMu.Lock()
	defer me.streamsMu.Unlock()
	
	for i, t := range stream.Tracks {
		if t.ID == track.ID {
			stream.Tracks = append(stream.Tracks[:i], stream.Tracks[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("track %s not found in stream", track.ID)
}

// Media capture (simplified implementations)
func (me *mediaEngine) GetUserMedia(constraints *MediaConstraints) (*MediaStream, error) {
	stream, err := me.CreateMediaStream("user-media")
	if err != nil {
		return nil, err
	}
	
	// Add audio track if requested
	if constraints.Audio != nil && constraints.Audio.Enabled {
		audioTrack, err := me.CreateTrack(MediaTypeAudio, "", "microphone")
		if err != nil {
			return nil, err
		}
		me.AddTrackToStream(stream, audioTrack)
	}
	
	// Add video track if requested
	if constraints.Video != nil && constraints.Video.Enabled {
		videoTrack, err := me.CreateTrack(MediaTypeVideo, "", "camera")
		if err != nil {
			return nil, err
		}
		me.AddTrackToStream(stream, videoTrack)
	}
	
	return stream, nil
}

func (me *mediaEngine) GetDisplayMedia(constraints *DisplayMediaConstraints) (*MediaStream, error) {
	stream, err := me.CreateMediaStream("display-media")
	if err != nil {
		return nil, err
	}
	
	// Add video track for screen capture
	if constraints.Video != nil && constraints.Video.Enabled {
		videoTrack, err := me.CreateTrack(MediaTypeVideo, "", "screen")
		if err != nil {
			return nil, err
		}
		me.AddTrackToStream(stream, videoTrack)
	}
	
	// Add audio track if requested
	if constraints.Audio != nil && constraints.Audio.Enabled {
		audioTrack, err := me.CreateTrack(MediaTypeAudio, "", "system-audio")
		if err != nil {
			return nil, err
		}
		me.AddTrackToStream(stream, audioTrack)
	}
	
	return stream, nil
}

// Media recording
func (me *mediaEngine) CreateRecorder(stream *MediaStream, options *RecorderOptions) (MediaRecorder, error) {
	return newMediaRecorder(stream, options)
}

// Helper functions
func formatParameters(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	
	var result string
	first := true
	for key, value := range params {
		if !first {
			result += ";"
		}
		result += key + "=" + value
		first = false
	}
	
	return result
}
