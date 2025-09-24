package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

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

	// Create webview
	w := webview.New(false)
	defer w.Destroy()
	
	w.SetTitle(config.Name)
	w.SetSize(config.Width, config.Height, webview.HintNone)
	w.Navigate(config.URL)
	w.Run()
}
