package mac

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"godesktop/platforms"
)

//go:embed runner/runner
var webviewBinary []byte

//go:embed icon.png
var defaultIcon []byte

// MacBuilder implements the Platform interface for macOS
type MacBuilder struct{}

// NewMacBuilder creates a new macOS platform builder
func NewMacBuilder() *MacBuilder {
	return &MacBuilder{}
}

// CreateApp creates a macOS .app bundle
func (m *MacBuilder) CreateApp(config platforms.AppConfig) error {
	return CreateAppBundle(config)
}

// GetFileExtension returns the file extension for macOS apps
func (m *MacBuilder) GetFileExtension() string {
	return ".app"
}

// GetPlatformName returns the platform name
func (m *MacBuilder) GetPlatformName() string {
	return "macOS"
}

// AppConfig is an alias for platforms.AppConfig for backward compatibility
type AppConfig = platforms.AppConfig

// AppTemplate holds the configuration that will be written to config.json
type AppTemplate struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// CreateAppBundle creates a macOS .app bundle with the given configuration
func CreateAppBundle(config platforms.AppConfig) error {
	appName := config.Name
	appBundle := appName + ".app"
	contentsDir := filepath.Join(appBundle, "Contents")
	macosDir := filepath.Join(contentsDir, "MacOS")
	resourcesDir := filepath.Join(contentsDir, "Resources")

	// Create directory structure
	if err := os.MkdirAll(macosDir, 0755); err != nil {
		return fmt.Errorf("failed to create MacOS directory: %v", err)
	}
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return fmt.Errorf("failed to create Resources directory: %v", err)
	}

	// Create the webview runner binary and config
	runnerPath := filepath.Join(macosDir, appName)
	if err := createRunner(runnerPath, config); err != nil {
		return fmt.Errorf("failed to create runner: %v", err)
	}

	// Handle icon
	if config.IconPath != "" {
		if err := processIcon(config.IconPath, resourcesDir); err != nil {
			return fmt.Errorf("failed to process icon: %v", err)
		}
	} else {
		// Use embedded default icon
		if err := processDefaultIcon(resourcesDir); err != nil {
			return fmt.Errorf("failed to process default icon: %v", err)
		}
	}

	// Create Info.plist
	if err := createInfoPlist(filepath.Join(contentsDir, "Info.plist"), config); err != nil {
		return fmt.Errorf("failed to create Info.plist: %v", err)
	}

	return nil
}

// createRunner creates the self-contained webview runner
func createRunner(runnerPath string, config platforms.AppConfig) error {
	// Write the embedded binary
	if err := os.WriteFile(runnerPath, webviewBinary, 0755); err != nil {
		return fmt.Errorf("failed to write runner binary: %v", err)
	}

	// Create config.json in the same directory
	configPath := filepath.Join(filepath.Dir(runnerPath), "config.json")
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

// processIcon handles icon conversion and copying
func processIcon(iconPath, resourcesDir string) error {
	iconExt := strings.ToLower(filepath.Ext(iconPath))
	targetIconPath := filepath.Join(resourcesDir, "icon.icns")

	switch iconExt {
	case ".icns":
		// Copy existing .icns file
		return copyFile(iconPath, targetIconPath)
	case ".png":
		// Convert PNG to .icns
		return convertPNGToICNS(iconPath, targetIconPath)
	default:
		return fmt.Errorf("unsupported icon format: %s (use .png or .icns)", iconExt)
	}
}

// processDefaultIcon handles the embedded default icon
func processDefaultIcon(resourcesDir string) error {
	// Create temporary file for the embedded icon
	tempIconPath := filepath.Join(os.TempDir(), "godesktop_default_icon.png")
	if err := os.WriteFile(tempIconPath, defaultIcon, 0644); err != nil {
		return fmt.Errorf("failed to write temporary icon: %v", err)
	}
	defer os.Remove(tempIconPath)

	// Convert the temporary PNG to .icns
	targetIconPath := filepath.Join(resourcesDir, "icon.icns")
	return convertPNGToICNS(tempIconPath, targetIconPath)
}

// convertPNGToICNS converts a PNG file to macOS .icns format
func convertPNGToICNS(pngPath, icnsPath string) error {
	iconsetDir := strings.TrimSuffix(icnsPath, ".icns") + ".iconset"

	// Create iconset directory
	if err := os.MkdirAll(iconsetDir, 0755); err != nil {
		return fmt.Errorf("failed to create iconset directory: %v", err)
	}
	defer os.RemoveAll(iconsetDir)

	// Generate all required icon sizes
	sizes := []struct {
		size int
		name string
	}{
		{16, "icon_16x16.png"},
		{32, "icon_16x16@2x.png"},
		{32, "icon_32x32.png"},
		{64, "icon_32x32@2x.png"},
		{128, "icon_128x128.png"},
		{256, "icon_128x128@2x.png"},
		{256, "icon_256x256.png"},
		{512, "icon_256x256@2x.png"},
		{512, "icon_512x512.png"},
		{1024, "icon_512x512@2x.png"},
	}

	for _, s := range sizes {
		outPath := filepath.Join(iconsetDir, s.name)
		cmd := exec.Command("sips", "-z", fmt.Sprintf("%d", s.size), fmt.Sprintf("%d", s.size), pngPath, "--out", outPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to resize icon to %dx%d: %v", s.size, s.size, err)
		}
	}

	// Convert iconset to .icns
	cmd := exec.Command("iconutil", "-c", "icns", iconsetDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create .icns file: %v", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// createInfoPlist creates the macOS Info.plist file
func createInfoPlist(plistPath string, config platforms.AppConfig) error {
	// Generate bundle ID from app name
	bundleID := fmt.Sprintf("com.godesktop.%s", strings.ToLower(strings.ReplaceAll(config.Name, " ", "")))

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>%s</string>
    <key>CFBundleIdentifier</key>
    <string>%s</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>%s</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2024</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
</dict>
</plist>`, config.Name, bundleID, config.Name)

	return os.WriteFile(plistPath, []byte(plistContent), 0644)
}
