package hls

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// PlaylistManager handles HLS playlist creation and management
type PlaylistManager struct {
	config *Config
}

// NewPlaylistManager creates a new playlist manager
func NewPlaylistManager(config *Config) *PlaylistManager {
	return &PlaylistManager{
		config: config,
	}
}

// GenerateMasterPlaylist creates the master playlist for adaptive bitrate streaming
func (pm *PlaylistManager) GenerateMasterPlaylist(qualityLevels []QualityLevel) error {
	if !pm.config.AdaptiveBitrate {
		return nil // No master playlist needed for single quality
	}

	masterPath := pm.config.GetMasterPlaylistPath()

	// Create master playlist content
	content := pm.buildMasterPlaylistContent(qualityLevels)

	// Write to file
	if err := os.WriteFile(masterPath, []byte(content), 0644); err != nil {
		return &HLSError{
			Message: "failed to write master playlist",
			Code:    ErrCodePlaylistError,
			Cause:   err,
		}
	}

	return nil
}

// buildMasterPlaylistContent builds the content of the master playlist
func (pm *PlaylistManager) buildMasterPlaylistContent(qualityLevels []QualityLevel) string {
	var builder strings.Builder

	// Header
	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:6\n")

	// Sort quality levels by bitrate (ascending)
	sortedLevels := make([]QualityLevel, len(qualityLevels))
	copy(sortedLevels, qualityLevels)
	sort.Slice(sortedLevels, func(i, j int) bool {
		return pm.getBitrateValue(sortedLevels[i].VideoBitrate) < pm.getBitrateValue(sortedLevels[j].VideoBitrate)
	})

	// Add stream info for each quality level
	for _, level := range sortedLevels {
		pm.addStreamInfo(&builder, level)
	}

	return builder.String()
}

// addStreamInfo adds stream information to the master playlist
func (pm *PlaylistManager) addStreamInfo(builder *strings.Builder, level QualityLevel) {
	// Calculate total bandwidth (video + audio)
	videoBitrate := pm.getBitrateValue(level.VideoBitrate)
	audioBitrate := pm.getBitrateValue(level.AudioBitrate)
	totalBandwidth := videoBitrate + audioBitrate

	// Build codecs string
	codecs := pm.buildCodecsString(level)

	// Stream info line
	builder.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d", totalBandwidth))

	// Add resolution if available
	if level.Resolution != "" {
		builder.WriteString(fmt.Sprintf(",RESOLUTION=%s", level.Resolution))
	}

	// Add codecs
	if codecs != "" {
		builder.WriteString(fmt.Sprintf(",CODECS=\"%s\"", codecs))
	}

	// Add frame rate
	if level.FrameRate > 0 {
		builder.WriteString(fmt.Sprintf(",FRAME-RATE=%.3f", level.FrameRate))
	}

	builder.WriteString("\n")

	// Playlist URL
	playlistURL := pm.getPlaylistURL(level.Name)
	builder.WriteString(fmt.Sprintf("%s\n", playlistURL))
}

// GeneratePlaylist creates a playlist for a specific quality level
func (pm *PlaylistManager) GeneratePlaylist(qualityLevel QualityLevel, segments []Segment) error {
	playlistPath := pm.config.GetPlaylistPath(qualityLevel.Name)

	// Create playlist content
	content := pm.buildPlaylistContent(qualityLevel, segments)

	// Write to file
	if err := os.WriteFile(playlistPath, []byte(content), 0644); err != nil {
		return &HLSError{
			Message: fmt.Sprintf("failed to write playlist for quality %s", qualityLevel.Name),
			Code:    ErrCodePlaylistError,
			Cause:   err,
		}
	}

	return nil
}

// buildPlaylistContent builds the content of a quality-specific playlist
func (pm *PlaylistManager) buildPlaylistContent(_ QualityLevel, segments []Segment) string {
	var builder strings.Builder

	// Header
	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:3\n")

	// Target duration (maximum segment duration)
	targetDuration := int(pm.config.SegmentOptions.Duration.Seconds())
	if len(segments) > 0 {
		maxDuration := 0.0
		for _, segment := range segments {
			if segment.Duration > maxDuration {
				maxDuration = segment.Duration
			}
		}
		if int(maxDuration) > targetDuration {
			targetDuration = int(maxDuration) + 1
		}
	}
	builder.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", targetDuration))

	// Media sequence
	builder.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", pm.config.SegmentOptions.StartNumber))

	// Playlist type
	if pm.config.PlaylistType != "" {
		builder.WriteString(fmt.Sprintf("#EXT-X-PLAYLIST-TYPE:%s\n", strings.ToUpper(string(pm.config.PlaylistType))))
	}

	// Add encryption info if enabled
	if pm.config.Encryption != nil && pm.config.Encryption.Method != EncryptionNone {
		pm.addEncryptionInfo(&builder)
	}

	// Add segments
	for _, segment := range segments {
		pm.addSegmentInfo(&builder, segment)
	}

	// End list for VOD
	if pm.config.PlaylistType == PlaylistVOD {
		builder.WriteString("#EXT-X-ENDLIST\n")
	}

	return builder.String()
}

// addEncryptionInfo adds encryption information to the playlist
func (pm *PlaylistManager) addEncryptionInfo(builder *strings.Builder) {
	enc := pm.config.Encryption
	builder.WriteString(fmt.Sprintf("#EXT-X-KEY:METHOD=%s", enc.Method))

	if enc.KeyURI != "" {
		builder.WriteString(fmt.Sprintf(",URI=\"%s\"", enc.KeyURI))
	}

	if enc.IV != "" {
		builder.WriteString(fmt.Sprintf(",IV=%s", enc.IV))
	}

	if enc.KeyFormat != "" {
		builder.WriteString(fmt.Sprintf(",KEYFORMAT=\"%s\"", enc.KeyFormat))
	}

	if enc.KeyFormatVersions != "" {
		builder.WriteString(fmt.Sprintf(",KEYFORMATVERSIONS=\"%s\"", enc.KeyFormatVersions))
	}

	builder.WriteString("\n")
}

// addSegmentInfo adds segment information to the playlist
func (pm *PlaylistManager) addSegmentInfo(builder *strings.Builder, segment Segment) {
	// Discontinuity tag
	if segment.Discontinuity {
		builder.WriteString("#EXT-X-DISCONTINUITY\n")
	}

	// Map tag for fMP4 initialization segment
	if segment.Map != "" {
		builder.WriteString(fmt.Sprintf("#EXT-X-MAP:URI=\"%s\"\n", segment.Map))
	}

	// Segment-specific key
	if segment.Key != nil && segment.Key.Method != EncryptionNone {
		builder.WriteString(fmt.Sprintf("#EXT-X-KEY:METHOD=%s", segment.Key.Method))
		if segment.Key.KeyURI != "" {
			builder.WriteString(fmt.Sprintf(",URI=\"%s\"", segment.Key.KeyURI))
		}
		if segment.Key.IV != "" {
			builder.WriteString(fmt.Sprintf(",IV=%s", segment.Key.IV))
		}
		builder.WriteString("\n")
	}

	// Segment duration and title
	builder.WriteString(fmt.Sprintf("#EXTINF:%.6f", segment.Duration))
	if segment.Title != "" {
		builder.WriteString(fmt.Sprintf(",%s", segment.Title))
	}
	builder.WriteString("\n")

	// Byte range
	if segment.ByteRange != "" {
		builder.WriteString(fmt.Sprintf("#EXT-X-BYTERANGE:%s\n", segment.ByteRange))
	}

	// Segment URI
	segmentURL := pm.getSegmentURL(segment.URI)
	builder.WriteString(fmt.Sprintf("%s\n", segmentURL))
}

// UpdatePlaylist updates an existing playlist (for live streaming)
func (pm *PlaylistManager) UpdatePlaylist(qualityLevel QualityLevel, newSegments []Segment) error {
	if pm.config.PlaylistType == PlaylistVOD {
		return nil // VOD playlists don't need updates
	}

	playlistPath := pm.config.GetPlaylistPath(qualityLevel.Name)

	// Read existing playlist if it exists
	var existingSegments []Segment
	if _, err := os.Stat(playlistPath); err == nil {
		// Parse existing playlist to get segments
		existingSegments, _ = pm.parseExistingPlaylist(playlistPath)
	}

	// Combine existing and new segments
	allSegments := append(existingSegments, newSegments...)

	// Apply list size limit
	if pm.config.SegmentOptions.ListSize > 0 && len(allSegments) > pm.config.SegmentOptions.ListSize {
		// Keep only the last N segments
		start := len(allSegments) - pm.config.SegmentOptions.ListSize
		allSegments = allSegments[start:]

		// Delete old segment files if configured
		if pm.config.SegmentOptions.DeleteOld {
			pm.deleteOldSegments(existingSegments[:start], qualityLevel.Name)
		}
	}

	return pm.GeneratePlaylist(qualityLevel, allSegments)
}

// Helper functions

// getBitrateValue extracts numeric bitrate value from string (e.g., "1000k" -> 1000000)
func (pm *PlaylistManager) getBitrateValue(bitrate string) int {
	if bitrate == "" {
		return 0
	}

	// Remove 'k' or 'M' suffix and convert
	bitrate = strings.ToLower(bitrate)
	multiplier := 1

	if strings.HasSuffix(bitrate, "k") {
		multiplier = 1000
		bitrate = strings.TrimSuffix(bitrate, "k")
	} else if strings.HasSuffix(bitrate, "m") {
		multiplier = 1000000
		bitrate = strings.TrimSuffix(bitrate, "m")
	}

	if value, err := strconv.Atoi(bitrate); err == nil {
		return value * multiplier
	}

	return 0
}

// buildCodecsString builds the codecs string for the master playlist
func (pm *PlaylistManager) buildCodecsString(level QualityLevel) string {
	var codecs []string

	// Video codec
	switch level.VideoCodec {
	case "libx264":
		if level.Profile != "" && level.Level != "" {
			codecs = append(codecs, fmt.Sprintf("avc1.%s%s", pm.getH264ProfileCode(level.Profile), pm.getH264LevelCode(level.Level)))
		} else {
			codecs = append(codecs, "avc1.42E01E") // Baseline profile, level 3.0
		}
	case "libx265":
		codecs = append(codecs, "hev1.1.6.L93.B0")
	}

	// Audio codec
	switch level.AudioCodec {
	case "aac":
		codecs = append(codecs, "mp4a.40.2")
	case "libmp3lame":
		codecs = append(codecs, "mp4a.40.34")
	}

	return strings.Join(codecs, ",")
}

// getH264ProfileCode returns the H.264 profile code
func (pm *PlaylistManager) getH264ProfileCode(profile string) string {
	switch strings.ToLower(profile) {
	case "baseline":
		return "42E0"
	case "main":
		return "4D40"
	case "high":
		return "6400"
	default:
		return "42E0"
	}
}

// getH264LevelCode returns the H.264 level code
func (pm *PlaylistManager) getH264LevelCode(level string) string {
	switch level {
	case "3.0":
		return "1E"
	case "3.1":
		return "1F"
	case "4.0":
		return "28"
	case "5.1":
		return "33"
	default:
		return "1E"
	}
}

// getPlaylistURL returns the URL for a playlist file
func (pm *PlaylistManager) getPlaylistURL(qualityName string) string {
	if pm.config.BaseURL != "" {
		if pm.config.AdaptiveBitrate {
			return fmt.Sprintf("%s/%s/%s", pm.config.BaseURL, qualityName, pm.config.PlaylistName)
		}
		return fmt.Sprintf("%s/%s_%s", pm.config.BaseURL, qualityName, pm.config.PlaylistName)
	}

	if pm.config.AdaptiveBitrate {
		return fmt.Sprintf("%s/%s", qualityName, pm.config.PlaylistName)
	}
	return fmt.Sprintf("%s_%s", qualityName, pm.config.PlaylistName)
}

// getSegmentURL returns the URL for a segment file
func (pm *PlaylistManager) getSegmentURL(segmentName string) string {
	if pm.config.BaseURL != "" {
		return fmt.Sprintf("%s/%s", pm.config.BaseURL, segmentName)
	}
	return segmentName
}

// parseExistingPlaylist parses an existing playlist file to extract segments
func (pm *PlaylistManager) parseExistingPlaylist(_ string) ([]Segment, error) {
	// This is a simplified implementation
	// In a real implementation, you would parse the M3U8 file properly
	return []Segment{}, nil
}

// deleteOldSegments deletes old segment files
func (pm *PlaylistManager) deleteOldSegments(segments []Segment, qualityName string) {
	outputDir := pm.config.GetQualityOutputDir(qualityName)
	for _, segment := range segments {
		segmentPath := filepath.Join(outputDir, segment.URI)
		os.Remove(segmentPath) // Ignore errors
	}
}
