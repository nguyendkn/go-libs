package hls

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// Converter is the main HLS converter implementation
type Converter struct {
	config           *Config
	ffmpeg           ffmpeg.FFmpeg
	segmentProcessor *SegmentProcessor
	playlistManager  *PlaylistManager
	mutex            sync.RWMutex
}

// NewConverter creates a new HLS converter
func NewConverter(config *Config) (*Converter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	converter := &Converter{
		config:           config,
		ffmpeg:           config.FFmpeg,
		segmentProcessor: NewSegmentProcessor(config),
		playlistManager:  NewPlaylistManager(config),
	}

	return converter, nil
}

// Convert converts a video file to HLS format
func (c *Converter) Convert(ctx context.Context, inputFile string) (*ConversionResult, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	startTime := time.Now()
	
	// Initialize result
	result := &ConversionResult{
		OutputDir:     c.config.OutputDir,
		Playlists:     make(map[string]string),
		Segments:      make(map[string][]string),
		QualityLevels: c.config.QualityLevels,
		Stats: ConversionStats{
			StartTime: startTime,
		},
	}

	// Validate input file
	if err := c.validateInput(inputFile); err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Setup output directory
	if err := c.config.SetupOutputDir(); err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Get input file info
	inputInfo, err := c.ffmpeg.GetMediaInfo(inputFile)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to get input file info: %v", err)
		return result, &HLSError{
			Message: "failed to analyze input file",
			Code:    ErrCodeInvalidInput,
			Cause:   err,
		}
	}

	result.Duration = inputInfo.Duration
	result.Stats.InputSize = inputInfo.Size

	// Process quality levels
	if c.config.Parallel && len(c.config.QualityLevels) > 1 {
		err = c.convertParallel(ctx, inputFile, result)
	} else {
		err = c.convertSequential(ctx, inputFile, result)
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result, err
	}

	// Generate master playlist for adaptive bitrate
	if c.config.AdaptiveBitrate && len(c.config.QualityLevels) > 1 {
		if err := c.playlistManager.GenerateMasterPlaylist(c.config.QualityLevels); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to generate master playlist: %v", err)
			return result, err
		}
		result.MasterPlaylist = c.config.GetMasterPlaylistPath()
	}

	// Calculate final stats
	endTime := time.Now()
	result.Stats.EndTime = endTime
	result.Stats.Duration = endTime.Sub(startTime)
	result.Stats.QualityCount = len(c.config.QualityLevels)

	// Calculate output size
	if outputSize, err := c.calculateOutputSize(); err == nil {
		result.Stats.OutputSize = outputSize
		if result.Stats.InputSize > 0 {
			result.Stats.CompressionRatio = float64(result.Stats.OutputSize) / float64(result.Stats.InputSize)
		}
	}

	result.Success = true
	return result, nil
}

// convertSequential processes quality levels one by one
func (c *Converter) convertSequential(ctx context.Context, inputFile string, result *ConversionResult) error {
	for i, qualityLevel := range c.config.QualityLevels {
		if err := c.processQualityLevel(ctx, inputFile, qualityLevel, i, len(c.config.QualityLevels), result); err != nil {
			return err
		}
	}
	return nil
}

// convertParallel processes quality levels in parallel
func (c *Converter) convertParallel(ctx context.Context, inputFile string, result *ConversionResult) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(c.config.QualityLevels))
	semaphore := make(chan struct{}, c.config.MaxConcurrent)

	for i, qualityLevel := range c.config.QualityLevels {
		wg.Add(1)
		go func(ql QualityLevel, index int) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := c.processQualityLevel(ctx, inputFile, ql, index, len(c.config.QualityLevels), result); err != nil {
				errChan <- err
			}
		}(qualityLevel, i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// processQualityLevel processes a single quality level
func (c *Converter) processQualityLevel(ctx context.Context, inputFile string, qualityLevel QualityLevel, index, total int, result *ConversionResult) error {
	// Create progress callback for this quality level
	progressCallback := func(progress ConversionProgress) {
		if c.config.ProgressCallback != nil {
			// Adjust progress for multiple quality levels
			overallProgress := (float64(index)/float64(total))*100 + (progress.Progress/float64(total))
			progress.Progress = overallProgress
			c.config.ProgressCallback(progress)
		}
	}

	// Process segments
	segments, err := c.segmentProcessor.ProcessSegmentsWithProgress(ctx, inputFile, qualityLevel, progressCallback)
	if err != nil {
		return err
	}

	// Validate segments
	outputDir := c.config.GetQualityOutputDir(qualityLevel.Name)
	if err := c.segmentProcessor.ValidateSegments(segments, outputDir); err != nil {
		return err
	}

	// Generate playlist
	if err := c.playlistManager.GeneratePlaylist(qualityLevel, segments); err != nil {
		return err
	}

	// Update result
	c.updateResult(result, qualityLevel, segments)

	return nil
}

// updateResult updates the conversion result with quality level data
func (c *Converter) updateResult(result *ConversionResult, qualityLevel QualityLevel, segments []Segment) {
	playlistPath := c.config.GetPlaylistPath(qualityLevel.Name)
	result.Playlists[qualityLevel.Name] = playlistPath

	segmentPaths := make([]string, len(segments))
	for i, segment := range segments {
		outputDir := c.config.GetQualityOutputDir(qualityLevel.Name)
		segmentPaths[i] = filepath.Join(outputDir, segment.URI)
	}
	result.Segments[qualityLevel.Name] = segmentPaths
	result.Stats.SegmentCount += len(segments)
}

// ConvertWithOptions converts video with custom options
func (c *Converter) ConvertWithOptions(ctx context.Context, inputFile string, options *ConversionOptions) (*ConversionResult, error) {
	// Create a copy of config with custom options
	config := c.config.Clone()
	
	if options != nil {
		// Apply custom options
		if options.QualityLevels != nil {
			config.QualityLevels = options.QualityLevels
		}
		if options.SegmentDuration > 0 {
			config.SegmentOptions.Duration = options.SegmentDuration
		}
		if options.OutputDir != "" {
			config.OutputDir = options.OutputDir
		}
		if options.Parallel != nil {
			config.Parallel = *options.Parallel
		}
		if options.MaxConcurrent > 0 {
			config.MaxConcurrent = options.MaxConcurrent
		}
	}

	// Create temporary converter with custom config
	tempConverter, err := NewConverter(config)
	if err != nil {
		return nil, err
	}

	return tempConverter.Convert(ctx, inputFile)
}

// ConvertLive converts video for live streaming
func (c *Converter) ConvertLive(ctx context.Context, inputSource string) error {
	if c.config.PlaylistType != PlaylistLive {
		return &HLSError{
			Message: "live conversion requires PlaylistLive type",
			Code:    ErrCodeInvalidConfig,
		}
	}

	// Setup output directory
	if err := c.config.SetupOutputDir(); err != nil {
		return err
	}

	// Start live conversion for each quality level
	if c.config.Parallel {
		return c.convertLiveParallel(ctx, inputSource)
	}
	
	return c.convertLiveSequential(ctx, inputSource)
}

// convertLiveSequential processes live stream sequentially
func (c *Converter) convertLiveSequential(ctx context.Context, inputSource string) error {
	for _, qualityLevel := range c.config.QualityLevels {
		go func(ql QualityLevel) {
			c.processLiveQualityLevel(ctx, inputSource, ql)
		}(qualityLevel)
	}
	
	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

// convertLiveParallel processes live stream in parallel
func (c *Converter) convertLiveParallel(ctx context.Context, inputSource string) error {
	var wg sync.WaitGroup
	
	for _, qualityLevel := range c.config.QualityLevels {
		wg.Add(1)
		go func(ql QualityLevel) {
			defer wg.Done()
			c.processLiveQualityLevel(ctx, inputSource, ql)
		}(qualityLevel)
	}
	
	wg.Wait()
	return nil
}

// processLiveQualityLevel processes a single quality level for live streaming
func (c *Converter) processLiveQualityLevel(ctx context.Context, inputSource string, qualityLevel QualityLevel) {
	// Implementation for live streaming would be more complex
	// This is a simplified version
	
	outputDir := c.config.GetQualityOutputDir(qualityLevel.Name)
	segmentPattern := c.config.GetSegmentPattern(qualityLevel.Name)
	
	builder := c.ffmpeg.New().
		Input(inputSource).
		VideoCodec(qualityLevel.VideoCodec).
		AudioCodec(qualityLevel.AudioCodec).
		Resolution(qualityLevel.Resolution).
		VideoBitrate(qualityLevel.VideoBitrate).
		AudioBitrate(qualityLevel.AudioBitrate)
	
	if qualityLevel.FrameRate > 0 {
		builder = builder.FrameRate(qualityLevel.FrameRate)
	}
	
	builder = c.segmentProcessor.addHLSOptions(builder, outputDir, segmentPattern, qualityLevel)
	
	// Execute live conversion
	builder.Execute(ctx)
}

// Helper methods

// validateInput validates the input file
func (c *Converter) validateInput(inputFile string) error {
	if inputFile == "" {
		return &HLSError{
			Message: "input file path is required",
			Code:    ErrCodeInvalidInput,
		}
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return &HLSError{
			Message: fmt.Sprintf("input file does not exist: %s", inputFile),
			Code:    ErrCodeInvalidInput,
		}
	}

	return nil
}

// calculateOutputSize calculates the total size of output files
func (c *Converter) calculateOutputSize() (int64, error) {
	var totalSize int64
	
	err := filepath.Walk(c.config.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	
	return totalSize, err
}

// ConversionOptions contains options for custom conversion
type ConversionOptions struct {
	QualityLevels   []QualityLevel
	SegmentDuration time.Duration
	OutputDir       string
	Parallel        *bool
	MaxConcurrent   int
	Encryption      *EncryptionOptions
}

// Cleanup removes temporary files and directories
func (c *Converter) Cleanup() error {
	if c.config.CleanupTemp && c.config.TempDir != "" {
		return os.RemoveAll(c.config.TempDir)
	}
	return nil
}

// GetConfig returns the current configuration
func (c *Converter) GetConfig() *Config {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.config.Clone()
}

// UpdateConfig updates the converter configuration
func (c *Converter) UpdateConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.config = config
	c.segmentProcessor = NewSegmentProcessor(config)
	c.playlistManager = NewPlaylistManager(config)
	
	return nil
}
