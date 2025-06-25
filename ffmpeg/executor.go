package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Executor handles FFmpeg command execution
type Executor struct {
	binaryPath string
}

// NewExecutor creates a new command executor
func NewExecutor(binaryPath string) *Executor {
	return &Executor{
		binaryPath: binaryPath,
	}
}

// Execute runs FFmpeg command with the given arguments and options
func (e *Executor) Execute(ctx context.Context, args []string, opts *ExecuteOptions) error {
	if opts == nil {
		opts = &ExecuteOptions{
			Context: ctx,
			Timeout: 300 * time.Second,
		}
	}

	// Create command with context
	cmdCtx := ctx
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		cmdCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(cmdCtx, e.binaryPath, args...)

	// Set working directory if specified
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	// Setup pipes for progress tracking
	if opts.ProgressHandler != nil {
		return e.executeWithProgress(cmd, opts)
	}

	// Simple execution without progress tracking
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &FFmpegError{
			Message: fmt.Sprintf("ffmpeg execution failed: %s\nOutput: %s", err.Error(), string(output)),
			Code:    getExitCode(err),
			Command: strings.Join(append([]string{e.binaryPath}, args...), " "),
		}
	}

	return nil
}

// executeWithProgress runs FFmpeg with progress tracking
func (e *Executor) executeWithProgress(cmd *exec.Cmd, opts *ExecuteOptions) error {
	// FFmpeg outputs progress to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Create channels for communication
	progressChan := make(chan ProgressInfo, 10)
	errorChan := make(chan error, 1)
	doneChan := make(chan struct{})

	// Start progress parser goroutine
	go e.parseProgress(stderr, progressChan, errorChan)

	// Start stdout reader goroutine (to prevent blocking)
	go func() {
		io.Copy(io.Discard, stdout)
	}()

	// Start progress handler goroutine
	go func() {
		defer close(doneChan)
		for progress := range progressChan {
			if opts.ProgressHandler != nil {
				opts.ProgressHandler(progress)
			}
		}
	}()

	// Wait for command completion
	cmdErr := cmd.Wait()
	close(progressChan)

	// Wait for progress handler to finish
	<-doneChan

	// Check for parsing errors
	select {
	case parseErr := <-errorChan:
		if opts.ErrorHandler != nil {
			opts.ErrorHandler(parseErr)
		}
	default:
	}

	// Handle command errors
	if cmdErr != nil {
		return &FFmpegError{
			Message: fmt.Sprintf("ffmpeg execution failed: %s", cmdErr.Error()),
			Code:    getExitCode(cmdErr),
			Command: strings.Join(cmd.Args, " "),
		}
	}

	return nil
}

// parseProgress parses FFmpeg progress output
func (e *Executor) parseProgress(reader io.Reader, progressChan chan<- ProgressInfo, errorChan chan<- error) {
	scanner := bufio.NewScanner(reader)

	// Regex patterns for parsing FFmpeg output
	frameRegex := regexp.MustCompile(`frame=\s*(\d+)`)
	fpsRegex := regexp.MustCompile(`fps=\s*([\d.]+)`)
	bitrateRegex := regexp.MustCompile(`bitrate=\s*([\d.]+\w*)`)
	timeRegex := regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
	sizeRegex := regexp.MustCompile(`size=\s*(\d+)kB`)
	speedRegex := regexp.MustCompile(`speed=\s*([\d.]+)x`)

	var totalDuration time.Duration
	durationRegex := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})\.(\d{2})`)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse total duration from initial output
		if totalDuration == 0 {
			if matches := durationRegex.FindStringSubmatch(line); len(matches) == 5 {
				hours, _ := strconv.Atoi(matches[1])
				minutes, _ := strconv.Atoi(matches[2])
				seconds, _ := strconv.Atoi(matches[3])
				centiseconds, _ := strconv.Atoi(matches[4])

				totalDuration = time.Duration(hours)*time.Hour +
					time.Duration(minutes)*time.Minute +
					time.Duration(seconds)*time.Second +
					time.Duration(centiseconds)*10*time.Millisecond
			}
		}

		// Parse progress information
		if strings.Contains(line, "frame=") && strings.Contains(line, "time=") {
			progress := ProgressInfo{}

			// Parse frame number
			if matches := frameRegex.FindStringSubmatch(line); len(matches) > 1 {
				if frame, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					progress.Frame = frame
				}
			}

			// Parse FPS
			if matches := fpsRegex.FindStringSubmatch(line); len(matches) > 1 {
				if fps, err := strconv.ParseFloat(matches[1], 64); err == nil {
					progress.FPS = fps
				}
			}

			// Parse bitrate
			if matches := bitrateRegex.FindStringSubmatch(line); len(matches) > 1 {
				progress.Bitrate = matches[1]
			}

			// Parse current time
			if matches := timeRegex.FindStringSubmatch(line); len(matches) == 5 {
				hours, _ := strconv.Atoi(matches[1])
				minutes, _ := strconv.Atoi(matches[2])
				seconds, _ := strconv.Atoi(matches[3])
				centiseconds, _ := strconv.Atoi(matches[4])

				currentTime := time.Duration(hours)*time.Hour +
					time.Duration(minutes)*time.Minute +
					time.Duration(seconds)*time.Second +
					time.Duration(centiseconds)*10*time.Millisecond

				progress.Time = currentTime

				// Calculate progress percentage
				if totalDuration > 0 {
					progress.Progress = float64(currentTime) / float64(totalDuration) * 100
					if progress.Progress > 100 {
						progress.Progress = 100
					}
				}
			}

			// Parse size
			if matches := sizeRegex.FindStringSubmatch(line); len(matches) > 1 {
				if size, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					progress.Size = size * 1024 // Convert kB to bytes
				}
			}

			// Parse speed
			if matches := speedRegex.FindStringSubmatch(line); len(matches) > 1 {
				if speed, err := strconv.ParseFloat(matches[1], 64); err == nil {
					progress.Speed = speed
				}
			}

			// Send progress update
			select {
			case progressChan <- progress:
			default:
				// Channel is full, skip this update
			}
		}
	}

	if err := scanner.Err(); err != nil {
		select {
		case errorChan <- fmt.Errorf("error reading ffmpeg output: %w", err):
		default:
		}
	}
}

// getExitCode extracts exit code from command error
func getExitCode(err error) int {
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return -1
}

// ProbeMediaFile extracts media information using ffprobe
func ProbeMediaFile(binaryPath, filePath string) (*MediaInfo, error) {
	// Use ffprobe if available, otherwise fallback to ffmpeg
	probePath := strings.Replace(binaryPath, "ffmpeg", "ffprobe", 1)

	// Check if ffprobe exists
	if !fileExists(probePath) {
		probePath = binaryPath
	}

	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	}

	cmd := exec.Command(probePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to probe media file: %w", err)
	}

	// Parse JSON output (simplified parsing for now)
	// In a real implementation, you would use proper JSON parsing
	info := &MediaInfo{}

	// Basic parsing - this should be replaced with proper JSON unmarshaling
	outputStr := string(output)

	// Extract duration using regex (simplified)
	durationRegex := regexp.MustCompile(`"duration":\s*"([\d.]+)"`)
	if matches := durationRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		if duration, err := strconv.ParseFloat(matches[1], 64); err == nil {
			info.Duration = time.Duration(duration * float64(time.Second))
		}
	}

	// Extract width and height
	widthRegex := regexp.MustCompile(`"width":\s*(\d+)`)
	heightRegex := regexp.MustCompile(`"height":\s*(\d+)`)

	if matches := widthRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		if width, err := strconv.Atoi(matches[1]); err == nil {
			info.Width = width
		}
	}

	if matches := heightRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		if height, err := strconv.Atoi(matches[1]); err == nil {
			info.Height = height
		}
	}

	// Extract codecs
	videoCodecRegex := regexp.MustCompile(`"codec_name":\s*"([^"]+)".*"codec_type":\s*"video"`)
	audioCodecRegex := regexp.MustCompile(`"codec_name":\s*"([^"]+)".*"codec_type":\s*"audio"`)

	if matches := videoCodecRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		info.VideoCodec = matches[1]
	}

	if matches := audioCodecRegex.FindStringSubmatch(outputStr); len(matches) > 1 {
		info.AudioCodec = matches[1]
	}

	return info, nil
}
