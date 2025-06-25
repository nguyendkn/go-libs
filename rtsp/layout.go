package rtsp

import (
	"fmt"
	"math"
	"strings"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// LayoutManager handles video layout and merging functionality
type LayoutManager struct {
	config *Config
	ffmpeg ffmpeg.FFmpeg
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(config *Config) *LayoutManager {
	return &LayoutManager{
		config: config,
		ffmpeg: config.FFmpeg,
	}
}

// CalculateLayout calculates the optimal layout for given streams
func (lm *LayoutManager) CalculateLayout(streamCount int) Layout {
	if lm.config.AutoLayout {
		return lm.config.AutoDetectLayout(streamCount)
	}
	return lm.config.Layout
}

// CalculateStreamPositions calculates positions for each stream in the layout
func (lm *LayoutManager) CalculateStreamPositions(layout Layout, streamNames []string) map[string]Position {
	positions := make(map[string]Position)
	
	// Calculate cell dimensions
	cellWidth := (layout.Width - (layout.Columns+1)*lm.config.LayoutPadding) / layout.Columns
	cellHeight := (layout.Height - (layout.Rows+1)*lm.config.LayoutPadding) / layout.Rows
	
	// Account for borders
	if lm.config.LayoutBorder > 0 {
		cellWidth -= 2 * lm.config.LayoutBorder
		cellHeight -= 2 * lm.config.LayoutBorder
	}
	
	streamIndex := 0
	for row := 0; row < layout.Rows && streamIndex < len(streamNames); row++ {
		for col := 0; col < layout.Columns && streamIndex < len(streamNames); col++ {
			streamName := streamNames[streamIndex]
			
			positions[streamName] = Position{
				Row:    row,
				Column: col,
				Width:  cellWidth,
				Height: cellHeight,
			}
			
			streamIndex++
		}
	}
	
	return positions
}

// GenerateFFmpegFilter generates FFmpeg filter for merging multiple streams
func (lm *LayoutManager) GenerateFFmpegFilter(layout Layout, streamNames []string) (string, error) {
	if len(streamNames) == 0 {
		return "", &RTSPError{
			Message: "no streams provided for layout",
			Code:    ErrCodeLayoutError,
		}
	}
	
	if len(streamNames) == 1 {
		// Single stream - just scale to output size
		return lm.generateSingleStreamFilter(layout), nil
	}
	
	positions := lm.CalculateStreamPositions(layout, streamNames)
	return lm.generateMultiStreamFilter(layout, streamNames, positions)
}

// generateSingleStreamFilter generates filter for single stream
func (lm *LayoutManager) generateSingleStreamFilter(layout Layout) string {
	return fmt.Sprintf("[0:v]scale=%d:%d[out]", layout.Width, layout.Height)
}

// generateMultiStreamFilter generates filter for multiple streams
func (lm *LayoutManager) generateMultiStreamFilter(layout Layout, streamNames []string, positions map[string]Position) (string, error) {
	var filterParts []string
	var overlayInputs []string
	
	// Create background
	bgFilter := fmt.Sprintf("color=%s:size=%dx%d[bg]", 
		lm.config.BackgroundColor, layout.Width, layout.Height)
	filterParts = append(filterParts, bgFilter)
	
	// Scale and position each stream
	for i, streamName := range streamNames {
		pos, exists := positions[streamName]
		if !exists {
			continue
		}
		
		// Calculate actual position with padding and borders
		x := pos.Column*(pos.Width+lm.config.LayoutPadding) + lm.config.LayoutPadding + lm.config.LayoutBorder
		y := pos.Row*(pos.Height+lm.config.LayoutPadding) + lm.config.LayoutPadding + lm.config.LayoutBorder
		
		// Scale stream to cell size
		scaleFilter := fmt.Sprintf("[%d:v]scale=%d:%d", i, pos.Width, pos.Height)
		
		// Add border if configured
		if lm.config.LayoutBorder > 0 {
			borderColor := lm.config.BackgroundColor
			if lm.config.Layout.BorderColor != "" {
				borderColor = lm.config.Layout.BorderColor
			}
			scaleFilter += fmt.Sprintf(",pad=%d:%d:%d:%d:%s", 
				pos.Width+2*lm.config.LayoutBorder,
				pos.Height+2*lm.config.LayoutBorder,
				lm.config.LayoutBorder,
				lm.config.LayoutBorder,
				borderColor)
		}
		
		scaleFilter += fmt.Sprintf("[s%d]", i)
		filterParts = append(filterParts, scaleFilter)
		
		// Prepare for overlay
		if i == 0 {
			overlayInputs = append(overlayInputs, "[bg][s0]")
		} else {
			overlayInputs = append(overlayInputs, fmt.Sprintf("[tmp%d][s%d]", i-1, i))
		}
		
		// Create overlay filter
		overlayFilter := fmt.Sprintf("overlay=%d:%d", x, y)
		if i < len(streamNames)-1 {
			overlayFilter += fmt.Sprintf("[tmp%d]", i)
		} else {
			overlayFilter += "[out]"
		}
		
		filterParts = append(filterParts, overlayInputs[i]+overlayFilter)
	}
	
	return strings.Join(filterParts, ";"), nil
}

// GenerateFFmpegArgs generates complete FFmpeg arguments for merging streams
func (lm *LayoutManager) GenerateFFmpegArgs(streamURLs []string, outputPath string, layout Layout) ([]string, error) {
	if len(streamURLs) == 0 {
		return nil, &RTSPError{
			Message: "no stream URLs provided",
			Code:    ErrCodeLayoutError,
		}
	}
	
	var args []string
	
	// Add input streams
	for _, url := range streamURLs {
		args = append(args, "-i", url)
		
		// Add RTSP-specific options
		args = append(args, "-rtsp_transport", string(lm.config.DefaultTransport))
		args = append(args, "-rtsp_flags", "prefer_tcp")
		args = append(args, "-stimeout", fmt.Sprintf("%d", int(lm.config.ConnectionTimeout.Seconds()*1000000))) // microseconds
	}
	
	// Generate filter
	streamNames := make([]string, len(streamURLs))
	for i := range streamURLs {
		streamNames[i] = fmt.Sprintf("stream_%d", i)
	}
	
	filter, err := lm.GenerateFFmpegFilter(layout, streamNames)
	if err != nil {
		return nil, err
	}
	
	// Add filter complex
	args = append(args, "-filter_complex", filter)
	
	// Add output options
	args = append(args, "-map", "[out]")
	
	// Add audio handling (mix all audio streams)
	if len(streamURLs) > 1 {
		audioFilter := lm.generateAudioMixFilter(len(streamURLs))
		args = append(args, "-filter_complex", audioFilter)
		args = append(args, "-map", "[aout]")
	} else {
		args = append(args, "-map", "0:a")
	}
	
	// Add codec settings
	args = append(args, "-c:v", string(lm.config.VideoCodec))
	args = append(args, "-c:a", string(lm.config.AudioCodec))
	
	// Add quality settings
	if lm.config.VideoBitrate != "" {
		args = append(args, "-b:v", lm.config.VideoBitrate)
	}
	if lm.config.AudioBitrate != "" {
		args = append(args, "-b:a", lm.config.AudioBitrate)
	}
	if lm.config.FrameRate > 0 {
		args = append(args, "-r", fmt.Sprintf("%.2f", lm.config.FrameRate))
	}
	
	// Add custom FFmpeg args
	if len(lm.config.FFmpegArgs) > 0 {
		args = append(args, lm.config.FFmpegArgs...)
	}
	
	// Add output path
	args = append(args, outputPath)
	
	return args, nil
}

// generateAudioMixFilter generates audio mixing filter for multiple streams
func (lm *LayoutManager) generateAudioMixFilter(streamCount int) string {
	if streamCount <= 1 {
		return "[0:a]acopy[aout]"
	}
	
	// Create audio mix filter
	var inputs []string
	for i := 0; i < streamCount; i++ {
		inputs = append(inputs, fmt.Sprintf("[%d:a]", i))
	}
	
	return fmt.Sprintf("%samix=inputs=%d[aout]", strings.Join(inputs, ""), streamCount)
}

// ValidateLayout validates a layout configuration
func (lm *LayoutManager) ValidateLayout(layout Layout, streamCount int) error {
	if layout.Rows <= 0 || layout.Columns <= 0 {
		return &RTSPError{
			Message: "layout rows and columns must be positive",
			Code:    ErrCodeLayoutError,
		}
	}
	
	if layout.Width <= 0 || layout.Height <= 0 {
		return &RTSPError{
			Message: "layout width and height must be positive",
			Code:    ErrCodeLayoutError,
		}
	}
	
	maxStreams := layout.Rows * layout.Columns
	if streamCount > maxStreams {
		return &RTSPError{
			Message: fmt.Sprintf("layout can accommodate maximum %d streams, but %d streams provided", maxStreams, streamCount),
			Code:    ErrCodeLayoutError,
		}
	}
	
	// Check if cell dimensions are reasonable
	cellWidth := layout.Width / layout.Columns
	cellHeight := layout.Height / layout.Rows
	
	if cellWidth < 64 || cellHeight < 64 {
		return &RTSPError{
			Message: "calculated cell dimensions are too small (minimum 64x64)",
			Code:    ErrCodeLayoutError,
		}
	}
	
	return nil
}

// GetOptimalLayout suggests an optimal layout for given stream count and constraints
func (lm *LayoutManager) GetOptimalLayout(streamCount int, maxWidth, maxHeight int) Layout {
	if streamCount <= 0 {
		return DefaultLayouts[LayoutSingle]
	}
	
	// Try to find the most square-like layout
	bestLayout := Layout{}
	bestRatio := math.MaxFloat64
	
	for rows := 1; rows <= streamCount; rows++ {
		cols := int(math.Ceil(float64(streamCount) / float64(rows)))
		
		if rows*cols >= streamCount {
			// Calculate aspect ratio difference from target
			layoutWidth := maxWidth
			layoutHeight := maxHeight
			
			cellWidth := layoutWidth / cols
			cellHeight := layoutHeight / rows
			
			// Skip if cells are too small
			if cellWidth < 64 || cellHeight < 64 {
				continue
			}
			
			// Calculate how "square" this layout is
			ratio := math.Abs(float64(cols)/float64(rows) - 1.0)
			
			if ratio < bestRatio {
				bestRatio = ratio
				bestLayout = Layout{
					Type:    LayoutCustom,
					Rows:    rows,
					Columns: cols,
					Width:   layoutWidth,
					Height:  layoutHeight,
				}
			}
		}
	}
	
	if bestLayout.Rows == 0 {
		// Fallback to default
		return DefaultLayouts[Layout2x2]
	}
	
	return bestLayout
}

// CreateCustomLayout creates a custom layout with specific parameters
func (lm *LayoutManager) CreateCustomLayout(rows, cols, width, height int) Layout {
	return Layout{
		Type:        LayoutCustom,
		Rows:        rows,
		Columns:     cols,
		Width:       width,
		Height:      height,
		Padding:     lm.config.LayoutPadding,
		Background:  lm.config.BackgroundColor,
		BorderWidth: lm.config.LayoutBorder,
		BorderColor: lm.config.Layout.BorderColor,
	}
}

// GetLayoutInfo returns information about a layout
func (lm *LayoutManager) GetLayoutInfo(layout Layout) map[string]interface{} {
	cellWidth := (layout.Width - (layout.Columns+1)*layout.Padding) / layout.Columns
	cellHeight := (layout.Height - (layout.Rows+1)*layout.Padding) / layout.Rows
	
	if layout.BorderWidth > 0 {
		cellWidth -= 2 * layout.BorderWidth
		cellHeight -= 2 * layout.BorderWidth
	}
	
	return map[string]interface{}{
		"type":              layout.Type,
		"rows":              layout.Rows,
		"columns":           layout.Columns,
		"total_cells":       layout.Rows * layout.Columns,
		"output_width":      layout.Width,
		"output_height":     layout.Height,
		"cell_width":        cellWidth,
		"cell_height":       cellHeight,
		"padding":           layout.Padding,
		"border_width":      layout.BorderWidth,
		"background_color":  layout.Background,
		"border_color":      layout.BorderColor,
	}
}

// PreviewLayout generates a preview description of the layout
func (lm *LayoutManager) PreviewLayout(layout Layout, streamNames []string) string {
	positions := lm.CalculateStreamPositions(layout, streamNames)
	
	var preview strings.Builder
	preview.WriteString(fmt.Sprintf("Layout: %dx%d (%s)\n", layout.Rows, layout.Columns, layout.Type))
	preview.WriteString(fmt.Sprintf("Output: %dx%d\n", layout.Width, layout.Height))
	preview.WriteString("Stream positions:\n")
	
	for _, streamName := range streamNames {
		if pos, exists := positions[streamName]; exists {
			preview.WriteString(fmt.Sprintf("  %s: Row %d, Col %d (%dx%d)\n", 
				streamName, pos.Row+1, pos.Column+1, pos.Width, pos.Height))
		}
	}
	
	return preview.String()
}
