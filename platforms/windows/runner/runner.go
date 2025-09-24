package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

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

	w.Navigate(config.URL)
	w.Run()
}
