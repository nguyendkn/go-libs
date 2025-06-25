package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Config holds FFmpeg configuration
type Config struct {
	BinaryPath string
	Timeout    int // seconds
	LogLevel   string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		BinaryPath: "",
		Timeout:    300, // 5 minutes
		LogLevel:   "error",
	}
}

// DetectFFmpegPath attempts to find FFmpeg binary in system PATH
func DetectFFmpegPath() (string, error) {
	// Common binary names
	binaryNames := []string{"ffmpeg"}
	if runtime.GOOS == "windows" {
		binaryNames = append(binaryNames, "ffmpeg.exe")
	}

	// Try to find in PATH
	for _, name := range binaryNames {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	// Try common installation paths
	commonPaths := getCommonFFmpegPaths()
	for _, basePath := range commonPaths {
		for _, name := range binaryNames {
			fullPath := filepath.Join(basePath, name)
			if fileExists(fullPath) {
				return fullPath, nil
			}
		}
	}

	return "", fmt.Errorf("ffmpeg binary not found in PATH or common locations")
}

// getCommonFFmpegPaths returns common FFmpeg installation paths by OS
func getCommonFFmpegPaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"C:\\ffmpeg\\bin",
			"C:\\Program Files\\ffmpeg\\bin",
			"C:\\Program Files (x86)\\ffmpeg\\bin",
			"C:\\tools\\ffmpeg\\bin",
			filepath.Join(os.Getenv("USERPROFILE"), "ffmpeg", "bin"),
		}
	case "darwin":
		return []string{
			"/usr/local/bin",
			"/opt/homebrew/bin",
			"/usr/bin",
			"/opt/local/bin",
			filepath.Join(os.Getenv("HOME"), "bin"),
		}
	case "linux":
		return []string{
			"/usr/bin",
			"/usr/local/bin",
			"/opt/ffmpeg/bin",
			"/snap/bin",
			filepath.Join(os.Getenv("HOME"), "bin"),
			filepath.Join(os.Getenv("HOME"), ".local", "bin"),
		}
	default:
		return []string{
			"/usr/bin",
			"/usr/local/bin",
		}
	}
}

// ValidateFFmpeg checks if FFmpeg binary is valid and accessible
func ValidateFFmpeg(binaryPath string) error {
	if binaryPath == "" {
		return fmt.Errorf("ffmpeg binary path is empty")
	}

	if !fileExists(binaryPath) {
		return fmt.Errorf("ffmpeg binary not found at: %s", binaryPath)
	}

	// Test if binary is executable and responds to version command
	cmd := exec.Command(binaryPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute ffmpeg: %w", err)
	}

	// Check if output contains expected FFmpeg signature
	outputStr := string(output)
	if !strings.Contains(strings.ToLower(outputStr), "ffmpeg") {
		return fmt.Errorf("invalid ffmpeg binary: unexpected version output")
	}

	return nil
}

// GetFFmpegVersion extracts version information from FFmpeg binary
func GetFFmpegVersion(binaryPath string) (string, error) {
	cmd := exec.Command(binaryPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ffmpeg version: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		// First line usually contains version info
		firstLine := lines[0]
		if strings.Contains(firstLine, "ffmpeg version") {
			return strings.TrimSpace(firstLine), nil
		}
	}

	return "", fmt.Errorf("could not parse ffmpeg version from output")
}

// GetSupportedFormats returns list of supported formats
func GetSupportedFormats(binaryPath string) ([]string, error) {
	cmd := exec.Command(binaryPath, "-formats")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get supported formats: %w", err)
	}

	var formats []string
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, " DE ") || strings.HasPrefix(line, "  E ") {
			// Extract format name (second column)
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				formats = append(formats, parts[1])
			}
		}
	}

	return formats, nil
}

// GetSupportedCodecs returns list of supported codecs
func GetSupportedCodecs(binaryPath string) (map[string][]string, error) {
	cmd := exec.Command(binaryPath, "-codecs")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get supported codecs: %w", err)
	}

	codecs := map[string][]string{
		"video": {},
		"audio": {},
	}

	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 8 {
			continue
		}

		// Check if line starts with codec flags
		if strings.HasPrefix(line, " D") || strings.HasPrefix(line, " E") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				codecName := parts[1]
				
				// Determine if it's video or audio codec
				if strings.Contains(line, "V") {
					codecs["video"] = append(codecs["video"], codecName)
				} else if strings.Contains(line, "A") {
					codecs["audio"] = append(codecs["audio"], codecName)
				}
			}
		}
	}

	return codecs, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// SetupConfig initializes FFmpeg configuration with auto-detection
func SetupConfig() (*Config, error) {
	config := DefaultConfig()
	
	// Try to detect FFmpeg binary
	binaryPath, err := DetectFFmpegPath()
	if err != nil {
		return nil, fmt.Errorf("failed to detect ffmpeg: %w", err)
	}
	
	config.BinaryPath = binaryPath
	
	// Validate the detected binary
	if err := ValidateFFmpeg(binaryPath); err != nil {
		return nil, fmt.Errorf("ffmpeg validation failed: %w", err)
	}
	
	return config, nil
}
