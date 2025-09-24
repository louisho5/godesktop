package windows

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"godesktop/platforms"
)

//go:embed runner/runner.exe
var webviewBinary []byte

// WindowsBuilder implements the Platform interface for Windows
type WindowsBuilder struct{}

// NewWindowsBuilder creates a new Windows platform builder
func NewWindowsBuilder() *WindowsBuilder {
	return &WindowsBuilder{}
}

// CreateApp creates a Windows application
func (w *WindowsBuilder) CreateApp(config platforms.AppConfig) error {
	return CreateWindowsApp(config)
}

// GetFileExtension returns the file extension for Windows apps
func (w *WindowsBuilder) GetFileExtension() string {
	return ".exe"
}

// GetPlatformName returns the platform name
func (w *WindowsBuilder) GetPlatformName() string {
	return "Windows"
}

// AppTemplate holds the configuration that will be written to config.json
type AppTemplate struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// CreateWindowsApp creates a Windows executable with the given configuration
func CreateWindowsApp(config platforms.AppConfig) error {
	appName := config.Name + ".exe"

	// Write the embedded binary
	if err := os.WriteFile(appName, webviewBinary, 0755); err != nil {
		return fmt.Errorf("failed to write Windows executable: %v", err)
	}

	// Create config.json in the same directory as the executable
	configPath := filepath.Join(filepath.Dir(appName), "config.json")
	appTemplate := AppTemplate{
		Name:   config.Name,
		Width:  config.Width,
		Height: config.Height,
		URL:    config.URL,
	}

	configData, err := json.Marshal(appTemplate)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}
