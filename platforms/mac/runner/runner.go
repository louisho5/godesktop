package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	webview "github.com/webview/webview_go"
)

type AppConfig struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func main() {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}

	// Look for config.json in the same directory as the executable
	configPath := filepath.Join(filepath.Dir(execPath), "config.json")

	// Read configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("Failed to read config:", err)
	}

	var config AppConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// Check if URL is a local file and start server if needed
	finalURL := config.URL
	if !strings.HasPrefix(config.URL, "http://") && !strings.HasPrefix(config.URL, "https://") {
		// It's a local file, start a local server
		// Get the app bundle root directory instead of the executable directory
		appDir := getAppBundleRoot(execPath)
		serverURL, err := startLocalServer(appDir)
		if err != nil {
			log.Fatal("Failed to start local server:", err)
		}
		finalURL = serverURL
		log.Printf("Local server started at: %s", serverURL)
	}

	// Create webview
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle(config.Name)
	w.SetSize(config.Width, config.Height, 0)
	w.Navigate(finalURL)
	w.Run()
}

func startLocalServer(dir string) (string, error) {
	// Start server on a random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", fmt.Errorf("failed to start server: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	serverURL := fmt.Sprintf("http://localhost:%d", port)

	// Set up file server
	fs := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve files from the application directory
		fs.ServeHTTP(w, r)
	})

	// Start server in background
	go func() {
		server := &http.Server{}
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	return serverURL, nil
}

// getAppBundleRoot returns the root directory of the macOS app bundle
// If the executable is in MyApp.app/Contents/MacOS/runner, this returns the directory containing MyApp.app
func getAppBundleRoot(execPath string) string {
	dir := filepath.Dir(execPath)
	
	// Check if we're inside a macOS app bundle (Contents/MacOS)
	if strings.HasSuffix(dir, "Contents/MacOS") {
		// Go up two levels: Contents/MacOS -> Contents -> MyApp.app -> parent directory
		appBundleDir := filepath.Dir(filepath.Dir(dir))
		return filepath.Dir(appBundleDir)
	}
	
	// If not in an app bundle, just return the executable directory
	return dir
}
