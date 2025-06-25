package hls

import (
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// HLSFormat represents HLS output formats
type HLSFormat string

const (
	FormatHLS     HLSFormat = "hls"
	FormatDASH    HLSFormat = "dash"
	FormatSmoothStreaming HLSFormat = "smooth"
)

// SegmentFormat represents segment file formats
type SegmentFormat string

const (
	SegmentTS   SegmentFormat = "ts"   // MPEG-TS (default)
	SegmentMP4  SegmentFormat = "mp4"  // fMP4
	SegmentWebM SegmentFormat = "webm" // WebM
)

// PlaylistType represents HLS playlist types
type PlaylistType string

const (
	PlaylistVOD   PlaylistType = "vod"   // Video on Demand
	PlaylistEvent PlaylistType = "event" // Live event
	PlaylistLive  PlaylistType = "live"  // Live streaming
)

// EncryptionMethod represents HLS encryption methods
type EncryptionMethod string

const (
	EncryptionNone   EncryptionMethod = "NONE"
	EncryptionAES128 EncryptionMethod = "AES-128"
	EncryptionSampleAES EncryptionMethod = "SAMPLE-AES"
)

// QualityLevel represents a quality level for adaptive streaming
type QualityLevel struct {
	Name         string             `json:"name"`
	Resolution   ffmpeg.Resolution  `json:"resolution"`
	VideoBitrate string             `json:"video_bitrate"`
	AudioBitrate string             `json:"audio_bitrate"`
	VideoCodec   ffmpeg.VideoCodec  `json:"video_codec"`
	AudioCodec   ffmpeg.AudioCodec  `json:"audio_codec"`
	FrameRate    float64            `json:"frame_rate"`
	Profile      string             `json:"profile"`
	Level        string             `json:"level"`
}

// Predefined quality levels
var (
	QualityLow = QualityLevel{
		Name:         "low",
		Resolution:   ffmpeg.Resolution480p,
		VideoBitrate: "800k",
		AudioBitrate: "96k",
		VideoCodec:   ffmpeg.VideoCodecH264,
		AudioCodec:   ffmpeg.AudioCodecAAC,
		FrameRate:    30,
		Profile:      "baseline",
		Level:        "3.0",
	}

	QualityMedium = QualityLevel{
		Name:         "medium",
		Resolution:   ffmpeg.Resolution720p,
		VideoBitrate: "2500k",
		AudioBitrate: "128k",
		VideoCodec:   ffmpeg.VideoCodecH264,
		AudioCodec:   ffmpeg.AudioCodecAAC,
		FrameRate:    30,
		Profile:      "main",
		Level:        "3.1",
	}

	QualityHigh = QualityLevel{
		Name:         "high",
		Resolution:   ffmpeg.Resolution1080p,
		VideoBitrate: "5000k",
		AudioBitrate: "192k",
		VideoCodec:   ffmpeg.VideoCodecH264,
		AudioCodec:   ffmpeg.AudioCodecAAC,
		FrameRate:    30,
		Profile:      "high",
		Level:        "4.0",
	}

	QualityUltra = QualityLevel{
		Name:         "ultra",
		Resolution:   ffmpeg.Resolution2160p,
		VideoBitrate: "15000k",
		AudioBitrate: "256k",
		VideoCodec:   ffmpeg.VideoCodecH264,
		AudioCodec:   ffmpeg.AudioCodecAAC,
		FrameRate:    30,
		Profile:      "high",
		Level:        "5.1",
	}
)

// DefaultQualityLevels provides common quality configurations
var DefaultQualityLevels = []QualityLevel{
	QualityLow,
	QualityMedium,
	QualityHigh,
}

// SegmentOptions contains options for HLS segmentation
type SegmentOptions struct {
	Duration     time.Duration `json:"duration"`      // Segment duration (default: 6s)
	ListSize     int           `json:"list_size"`     // Number of segments in playlist (default: 0 = all)
	Format       SegmentFormat `json:"format"`        // Segment format (default: ts)
	Pattern      string        `json:"pattern"`       // Segment filename pattern
	StartNumber  int           `json:"start_number"`  // Starting segment number
	DeleteOld    bool          `json:"delete_old"`    // Delete old segments
	InitSegment  string        `json:"init_segment"`  // Initialization segment name (for fMP4)
}

// EncryptionOptions contains options for HLS encryption
type EncryptionOptions struct {
	Method     EncryptionMethod `json:"method"`
	KeyFile    string           `json:"key_file"`
	KeyURI     string           `json:"key_uri"`
	KeyFormat  string           `json:"key_format"`
	KeyFormatVersions string    `json:"key_format_versions"`
	IV         string           `json:"iv"`
}

// ConversionProgress represents HLS conversion progress
type ConversionProgress struct {
	Stage        string        `json:"stage"`         // Current stage (analyzing, converting, etc.)
	QualityLevel string        `json:"quality_level"` // Current quality being processed
	Segment      int           `json:"segment"`       // Current segment number
	TotalSegments int          `json:"total_segments"`
	Progress     float64       `json:"progress"`      // Overall progress (0-100)
	Speed        float64       `json:"speed"`         // Processing speed
	ETA          time.Duration `json:"eta"`           // Estimated time remaining
	FFmpegProgress ffmpeg.ProgressInfo `json:"ffmpeg_progress"`
}

// ConversionResult contains the result of HLS conversion
type ConversionResult struct {
	Success       bool                    `json:"success"`
	OutputDir     string                  `json:"output_dir"`
	MasterPlaylist string                 `json:"master_playlist"`
	Playlists     map[string]string       `json:"playlists"`     // quality -> playlist path
	Segments      map[string][]string     `json:"segments"`      // quality -> segment paths
	Duration      time.Duration           `json:"duration"`
	QualityLevels []QualityLevel          `json:"quality_levels"`
	Error         string                  `json:"error,omitempty"`
	Stats         ConversionStats         `json:"stats"`
}

// ConversionStats contains statistics about the conversion
type ConversionStats struct {
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Duration      time.Duration `json:"duration"`
	InputSize     int64         `json:"input_size"`
	OutputSize    int64         `json:"output_size"`
	CompressionRatio float64    `json:"compression_ratio"`
	SegmentCount  int           `json:"segment_count"`
	QualityCount  int           `json:"quality_count"`
}

// HLSError represents HLS-specific errors
type HLSError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Stage   string `json:"stage"`
	Cause   error  `json:"cause,omitempty"`
}

func (e *HLSError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *HLSError) Unwrap() error {
	return e.Cause
}

// Error codes
const (
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodeFFmpegError     = "FFMPEG_ERROR"
	ErrCodeFileSystem      = "FILESYSTEM_ERROR"
	ErrCodeInvalidConfig   = "INVALID_CONFIG"
	ErrCodeConversionFailed = "CONVERSION_FAILED"
	ErrCodePlaylistError   = "PLAYLIST_ERROR"
	ErrCodeSegmentError    = "SEGMENT_ERROR"
)

// Playlist represents an HLS playlist
type Playlist struct {
	Type        PlaylistType `json:"type"`
	Version     int          `json:"version"`
	TargetDuration int       `json:"target_duration"`
	MediaSequence int        `json:"media_sequence"`
	Segments    []Segment    `json:"segments"`
	EndList     bool         `json:"end_list"`
	Encryption  *EncryptionOptions `json:"encryption,omitempty"`
}

// Segment represents an HLS segment
type Segment struct {
	URI        string        `json:"uri"`
	Duration   float64       `json:"duration"`
	Title      string        `json:"title,omitempty"`
	ByteRange  string        `json:"byte_range,omitempty"`
	Discontinuity bool       `json:"discontinuity,omitempty"`
	Key        *EncryptionOptions `json:"key,omitempty"`
	Map        string        `json:"map,omitempty"` // For fMP4 initialization segment
}

// MasterPlaylist represents an HLS master playlist
type MasterPlaylist struct {
	Version   int                    `json:"version"`
	Streams   []StreamInfo           `json:"streams"`
	Media     []MediaInfo            `json:"media"`
	SessionData []SessionDataInfo    `json:"session_data,omitempty"`
}

// StreamInfo represents a stream in master playlist
type StreamInfo struct {
	URI        string `json:"uri"`
	Bandwidth  int    `json:"bandwidth"`
	Resolution string `json:"resolution,omitempty"`
	Codecs     string `json:"codecs,omitempty"`
	FrameRate  float64 `json:"frame_rate,omitempty"`
	Audio      string `json:"audio,omitempty"`
	Video      string `json:"video,omitempty"`
	Subtitles  string `json:"subtitles,omitempty"`
}

// MediaInfo represents media in master playlist
type MediaInfo struct {
	Type       string `json:"type"`
	GroupID    string `json:"group_id"`
	Name       string `json:"name"`
	Default    bool   `json:"default"`
	AutoSelect bool   `json:"autoselect"`
	Language   string `json:"language,omitempty"`
	URI        string `json:"uri,omitempty"`
}

// SessionDataInfo represents session data in master playlist
type SessionDataInfo struct {
	DataID   string `json:"data_id"`
	Value    string `json:"value,omitempty"`
	URI      string `json:"uri,omitempty"`
	Language string `json:"language,omitempty"`
}
