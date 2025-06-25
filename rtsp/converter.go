package rtsp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
	"github.com/nguyendkn/go-libs/hls"
)

// Converter handles RTSP to HLS conversion
type Converter struct {
	config            *Config
	streamManager     *StreamManager
	layoutManager     *LayoutManager
	hlsConverter      hls.HLS
	ffmpeg            ffmpeg.FFmpeg
	mutex             sync.RWMutex
	activeConversions map[string]*ConversionContext
}

// ConversionContext represents an active conversion
type ConversionContext struct {
	ID          string
	StreamURLs  []string
	StreamNames []string
	OutputDir   string
	Layout      Layout
	Mode        StreamingMode
	StartTime   time.Time
	Context     context.Context
	Cancel      context.CancelFunc
	Progress    ConversionProgress
	Result      *ConversionResult
	Error       error
}

// NewConverter creates a new RTSP to HLS converter
func NewConverter(config *Config) (*Converter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	streamManager, err := NewStreamManager(config)
	if err != nil {
		return nil, err
	}

	layoutManager := NewLayoutManager(config)

	return &Converter{
		config:            config,
		streamManager:     streamManager,
		layoutManager:     layoutManager,
		hlsConverter:      config.HLS,
		ffmpeg:            config.FFmpeg,
		activeConversions: make(map[string]*ConversionContext),
	}, nil
}

// ConvertSingle converts a single RTSP stream to HLS
func (c *Converter) ConvertSingle(ctx context.Context, streamURL, outputDir string) (*ConversionResult, error) {
	return c.Convert(ctx, []string{streamURL}, outputDir, ModeSeparate)
}

// ConvertMultiple converts multiple RTSP streams to separate HLS streams
func (c *Converter) ConvertMultiple(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	return c.Convert(ctx, streamURLs, outputDir, ModeSeparate)
}

// ConvertMerged converts multiple RTSP streams to a single merged HLS stream
func (c *Converter) ConvertMerged(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	return c.Convert(ctx, streamURLs, outputDir, ModeMerged)
}

// Convert converts RTSP streams to HLS with specified mode
func (c *Converter) Convert(ctx context.Context, streamURLs []string, outputDir string, mode StreamingMode) (*ConversionResult, error) {
	if len(streamURLs) == 0 {
		return nil, &RTSPError{
			Message: "no stream URLs provided",
			Code:    ErrCodeInvalidConfig,
		}
	}

	// Validate stream URLs
	for _, streamURL := range streamURLs {
		if err := c.config.ValidateStreamURL(streamURL); err != nil {
			return nil, err
		}
	}

	// Setup output directory
	config := c.config.Clone()
	config.OutputDir = outputDir
	config.StreamingMode = mode
	if err := config.SetupOutputDir(); err != nil {
		return nil, err
	}

	// Create conversion context
	conversionID := fmt.Sprintf("conv_%d", time.Now().UnixNano())
	convCtx, cancel := context.WithCancel(ctx)

	streamNames := make([]string, len(streamURLs))
	for i := range streamURLs {
		streamNames[i] = fmt.Sprintf("stream_%d", i+1)
	}

	layout := c.layoutManager.CalculateLayout(len(streamURLs))
	if err := c.layoutManager.ValidateLayout(layout, len(streamURLs)); err != nil {
		cancel()
		return nil, err
	}

	conversionCtx := &ConversionContext{
		ID:          conversionID,
		StreamURLs:  streamURLs,
		StreamNames: streamNames,
		OutputDir:   outputDir,
		Layout:      layout,
		Mode:        mode,
		StartTime:   time.Now(),
		Context:     convCtx,
		Cancel:      cancel,
	}

	// Register conversion
	c.mutex.Lock()
	c.activeConversions[conversionID] = conversionCtx
	c.mutex.Unlock()

	// Cleanup on completion
	defer func() {
		c.mutex.Lock()
		delete(c.activeConversions, conversionID)
		c.mutex.Unlock()
		cancel()
	}()

	// Perform conversion based on mode
	switch mode {
	case ModeSeparate:
		return c.convertSeparate(conversionCtx)
	case ModeMerged:
		return c.convertMerged(conversionCtx)
	case ModeBoth:
		return c.convertBoth(conversionCtx)
	default:
		return nil, &RTSPError{
			Message: fmt.Sprintf("unsupported streaming mode: %s", mode),
			Code:    ErrCodeInvalidConfig,
		}
	}
}

// convertSeparate converts each stream to separate HLS
func (c *Converter) convertSeparate(convCtx *ConversionContext) (*ConversionResult, error) {
	result := &ConversionResult{
		Success:   true,
		Mode:      ModeSeparate,
		Layout:    &convCtx.Layout,
		OutputDir: convCtx.OutputDir,
		Streams:   make(map[string]StreamResult),
		StartTime: convCtx.StartTime,
		Stats:     ConversionStats{TotalStreams: len(convCtx.StreamURLs)},
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	errors := make([]error, 0)

	// Convert each stream separately
	for i, streamURL := range convCtx.StreamURLs {
		wg.Add(1)
		go func(index int, url string) {
			defer wg.Done()

			streamName := convCtx.StreamNames[index]
			streamOutputDir := c.config.GetStreamOutputDir(streamName)

			// Create stream output directory
			if err := os.MkdirAll(streamOutputDir, 0755); err != nil {
				mutex.Lock()
				errors = append(errors, err)
				mutex.Unlock()
				return
			}

			// Convert using HLS converter
			hlsResult, err := c.convertStreamToHLS(convCtx.Context, url, streamName, streamOutputDir)

			mutex.Lock()
			if err != nil {
				errors = append(errors, err)
				result.Streams[streamName] = StreamResult{
					StreamName: streamName,
					StreamURL:  url,
					Success:    false,
					Error:      err.Error(),
				}
				result.Stats.FailedStreams++
			} else {
				result.Streams[streamName] = StreamResult{
					StreamName:   streamName,
					StreamURL:    url,
					Success:      true,
					OutputDir:    streamOutputDir,
					PlaylistPath: hlsResult.MasterPlaylist,
					Duration:     hlsResult.Duration,
					SegmentCount: hlsResult.Stats.SegmentCount,
					HLSResult:    hlsResult,
				}
				result.Stats.SuccessfulStreams++
			}
			mutex.Unlock()

			// Report progress
			if c.config.ProgressCallback != nil {
				progress := ConversionProgress{
					StreamURL:   url,
					StreamName:  streamName,
					Status:      StatusStreaming,
					Progress:    float64(result.Stats.SuccessfulStreams+result.Stats.FailedStreams) / float64(result.Stats.TotalStreams) * 100,
					HLSProgress: nil, // Would be populated from HLS conversion
				}
				c.config.ProgressCallback(progress)
			}
		}(i, streamURL)
	}

	wg.Wait()

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if len(errors) > 0 {
		result.Success = false
		result.Error = fmt.Sprintf("failed to convert %d streams", len(errors))
		return result, errors[0]
	}

	return result, nil
}

// convertMerged converts multiple streams to a single merged HLS
func (c *Converter) convertMerged(convCtx *ConversionContext) (*ConversionResult, error) {
	result := &ConversionResult{
		Success:   true,
		Mode:      ModeMerged,
		Layout:    &convCtx.Layout,
		OutputDir: convCtx.OutputDir,
		Streams:   make(map[string]StreamResult),
		StartTime: convCtx.StartTime,
		Stats:     ConversionStats{TotalStreams: len(convCtx.StreamURLs)},
	}

	mergedOutputDir := c.config.GetMergedOutputDir()
	if err := os.MkdirAll(mergedOutputDir, 0755); err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Generate FFmpeg arguments for merging
	outputPath := filepath.Join(mergedOutputDir, "playlist.m3u8")
	args, err := c.layoutManager.GenerateFFmpegArgs(convCtx.StreamURLs, outputPath, convCtx.Layout)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Execute FFmpeg command
	if err := c.ffmpeg.Execute(convCtx.Context, args, nil); err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, &RTSPError{
			Message: "failed to merge streams",
			Code:    ErrCodeConversionFailed,
			Cause:   err,
		}
	}

	// Create merged stream result
	mergedResult := StreamResult{
		StreamName:   "merged",
		StreamURL:    fmt.Sprintf("merged_%d_streams", len(convCtx.StreamURLs)),
		Success:      true,
		OutputDir:    mergedOutputDir,
		PlaylistPath: outputPath,
	}

	result.MergedStream = &mergedResult
	result.Stats.SuccessfulStreams = 1
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// convertBoth converts streams both separately and merged
func (c *Converter) convertBoth(convCtx *ConversionContext) (*ConversionResult, error) {
	// First convert separately
	separateResult, err := c.convertSeparate(convCtx)
	if err != nil {
		return separateResult, err
	}

	// Then convert merged
	mergedResult, err := c.convertMerged(convCtx)
	if err != nil {
		// Return separate result even if merged fails
		separateResult.Error = fmt.Sprintf("merged conversion failed: %v", err)
		return separateResult, nil
	}

	// Combine results
	result := &ConversionResult{
		Success:        true,
		Mode:           ModeBoth,
		Layout:         &convCtx.Layout,
		OutputDir:      convCtx.OutputDir,
		MasterPlaylist: mergedResult.MasterPlaylist,
		Streams:        separateResult.Streams,
		MergedStream:   mergedResult.MergedStream,
		Duration:       separateResult.Duration + mergedResult.Duration,
		StartTime:      convCtx.StartTime,
		EndTime:        time.Now(),
		Stats: ConversionStats{
			TotalStreams:      separateResult.Stats.TotalStreams,
			SuccessfulStreams: separateResult.Stats.SuccessfulStreams,
			FailedStreams:     separateResult.Stats.FailedStreams,
		},
	}

	return result, nil
}

// convertStreamToHLS converts a single RTSP stream to HLS using the HLS package
func (c *Converter) convertStreamToHLS(ctx context.Context, streamURL, streamName, outputDir string) (*hls.ConversionResult, error) {
	// Create HLS configuration
	hlsConfig := c.config.GetHLSConfig()
	hlsConfig.OutputDir = outputDir

	// Create HLS converter
	hlsConverter, err := hls.NewWithConfig(hlsConfig)
	if err != nil {
		return nil, &RTSPError{
			Message: "failed to create HLS converter",
			Code:    ErrCodeHLSError,
			Cause:   err,
		}
	}
	defer hlsConverter.Cleanup()

	// For RTSP streams, we need to use FFmpeg to capture and convert
	// This is a simplified approach - in production, you might want to
	// handle live streams differently

	// Create a temporary file to capture the stream
	tempFile := filepath.Join(outputDir, "temp_capture.mp4")
	defer os.Remove(tempFile)

	// Capture stream for a duration (this is for demo purposes)
	// In a real implementation, you'd handle live streaming differently
	captureArgs := []string{
		"-i", streamURL,
		"-t", "30", // Capture 30 seconds for demo
		"-c", "copy",
		tempFile,
	}

	if err := c.ffmpeg.Execute(ctx, captureArgs, nil); err != nil {
		return nil, &RTSPError{
			Message: "failed to capture RTSP stream",
			Code:    ErrCodeStreamError,
			Cause:   err,
		}
	}

	// Convert captured file to HLS
	return hlsConverter.Convert(ctx, tempFile)
}

// GetActiveConversions returns information about active conversions
func (c *Converter) GetActiveConversions() map[string]ConversionProgress {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	progress := make(map[string]ConversionProgress)
	for id, ctx := range c.activeConversions {
		progress[id] = ctx.Progress
	}
	return progress
}

// StopConversion stops an active conversion
func (c *Converter) StopConversion(conversionID string) error {
	c.mutex.RLock()
	ctx, exists := c.activeConversions[conversionID]
	c.mutex.RUnlock()

	if !exists {
		return &RTSPError{
			Message: fmt.Sprintf("conversion %s not found", conversionID),
			Code:    ErrCodeInvalidConfig,
		}
	}

	ctx.Cancel()
	return nil
}

// StopAllConversions stops all active conversions
func (c *Converter) StopAllConversions() {
	c.mutex.RLock()
	contexts := make([]*ConversionContext, 0, len(c.activeConversions))
	for _, ctx := range c.activeConversions {
		contexts = append(contexts, ctx)
	}
	c.mutex.RUnlock()

	for _, ctx := range contexts {
		ctx.Cancel()
	}
}

// Close closes the converter and cleans up resources
func (c *Converter) Close() error {
	c.StopAllConversions()
	return c.streamManager.Close()
}
