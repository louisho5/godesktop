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

	"github.com/jchv/go-webview2"
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
		serverURL, err := startLocalServer(filepath.Dir(execPath))
		if err != nil {
			log.Fatal("Failed to start local server:", err)
		}
		finalURL = serverURL
		log.Printf("Local server started at: %s", serverURL)
	}

	// Create webview
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  config.Name,
			Width:  uint(config.Width),
			Height: uint(config.Height),
			IconId: 0,
			Center: true,
		},
	})
	defer w.Destroy()

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
