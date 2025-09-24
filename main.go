package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"godesktop/platforms"
	"godesktop/platforms/mac"
	"godesktop/platforms/windows"
)

func main() {
	var config platforms.AppConfig

	// Define command-line flags
	flag.StringVar(&config.Name, "name", "GoDesktopApp", "Application name")
	flag.StringVar(&config.URL, "url", "", "The target url to navigate to")
	flag.IntVar(&config.Width, "width", 1200, "Window width (Optional)")
	flag.IntVar(&config.Height, "height", 800, "Window height (Optional)")
	flag.StringVar(&config.IconPath, "icon", "", "Path to PNG icon file (Optional & Mac Only)")

	flag.Parse()

	fmt.Println("GoDesktop CLI - Create native lightweight desktop apps from web content")
	fmt.Println("Help: ./godesktop -help \n")
	fmt.Println("Example: ./godesktop -name 'GitHub' -url 'https://github.com' \n")

	// Validate input
	if config.URL == "" {
		fmt.Println("Error: -url must be provided\n")
		fmt.Println("Enter a URL to navigate to: ")
		fmt.Scanln(&config.URL)
		config.URL = strings.TrimSpace(config.URL)
		if !strings.Contains(config.URL, "://") {
			config.URL = "http://" + config.URL
		}
	}

	// if config.Name == "" {
	// 	fmt.Println("Error: -name must be provided\n")
	// 	fmt.Println("Enter a name for your app: ")
	// 	fmt.Scanln(&config.Name)
	// }

	// Select platform builder based on runtime OS
	var builder platforms.Platform
	switch runtime.GOOS {
	case "darwin":
		builder = mac.NewMacBuilder()
	case "windows":
		builder = windows.NewWindowsBuilder()
	default:
		fmt.Printf("Error: Platform %s is not supported yet\n", runtime.GOOS)
		fmt.Println("Currently supported platforms: macOS (darwin), Windows (windows)")
		os.Exit(1)
	}

	// Create the app using the platform builder
	err := builder.CreateApp(config)
	if err != nil {
		fmt.Printf("Error creating app: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s app created: %s%s\n", builder.GetPlatformName(), config.Name, builder.GetFileExtension())

	// Platform-specific run instructions
	switch runtime.GOOS {
	case "darwin":
		fmt.Printf("   Run with: open %s%s\n", config.Name, builder.GetFileExtension())
	case "windows":
		fmt.Printf("   Run with: .\\%s%s\n", config.Name, builder.GetFileExtension())
	default:
		fmt.Printf("   Run the created application\n")
	}
}
