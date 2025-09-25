# GoDesktop CLI

#### A CLI tool to create native lightweight desktop apps from a URL.

![GoDesktop](platforms/mac/icon.png)

✅ Single portable binary — no dependencies, no runtime required.

✅ Alternatives to [nativefier](https://github.com/nativefier/nativefier) but faster and smaller.

✅ Drop in a URL, get a desktop app.

```bash
./godesktop -url "https://github.com"
```

## Quick Start

### Create Your First App

**MacOS (Darwin)**

Download the latest "godesktop" release from the [releases page](https://github.com/louisho5/godesktop/releases).

```bash
# Changer permissions to make it executable
chmod +x godesktop

# Create app from URL
./godesktop -name "GitHub" -url "https://github.com"

# Create app from static HTML file
./godesktop -name "My HTML" -url "index.html"

# Create app with custom icon
./godesktop -name "My App" -url "https://example.com" -icon "icon.png" -width 900 -height 700
```

**Windows 10/11**

Download the latest "godesktop.exe" release from the [releases page](https://github.com/louisho5/godesktop/releases).

```bash
# Create app from URL
.\godesktop.exe -name "GitHub" -url "https://github.com"

# Create app from static HTML file
.\godesktop.exe -name "My HTML" -url "index.html"

# Create app with custom icon
.\godesktop.exe -name "My App" -url "https://example.com" -width 900 -height 700
```

## Usage

```bash
godesktop [flags]
```

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-name` | string | "GoDesktop" | Application name |
| `-width` | int | 1200 | Window width in pixels |
| `-height` | int | 800 | Window height in pixels |
| `-url` | string | - | URL to navigate to (required). Use "index.html" to serve a static site. |
| `-icon` | string | - | Path to icon file (.png) |

### Icon Support

*GoDesktop only handles icon conversion for macOS:*

- **PNG files**: Accepts .png files with all required sizes
- **Recommended size**: 512x512 or 1024x1024 pixels for best quality


## Development

### Build the CLI

```bash
# Clone and build in one command
git clone <repo-url>
cd godesktop

======== MacOS

# Step 1: Rebuild the MacOS runner (Optional)
go build -o platforms/mac/runner/runner platforms/mac/runner/runner.go

# Step 2: Build the CLI for macOS
go build -o ./godesktop main.go

======== Windows

# Step 1A (Select one): Rebuild the Windows runner (Optional)
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags="-H=windowsgui" -o platforms/windows/runner/runner.exe platforms/windows/runner/runner.go
# Step 1B (Select one): Run the build.bat script (Optional)
cd ".\platforms\windows\runner"; .\build.bat; cd ../../../

# Step 2: Build the CLI for Windows
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o ./godesktop.exe main.go
```

### Dependencies

GoDesktop is built with Go and uses the following excellent libraries:

#### Core Libraries
- **[webview/webview_go](https://github.com/webview/webview_go)** - Cross-platform webview library for macOS
- **[jchv/go-webview2](https://github.com/jchv/go-webview2)** - WebView2 bindings for Windows

#### System Requirements
- **Go 1.24+** for building
- **macOS**: Uses native WebKit framework (xcode command line tools required)
- **Windows**: Requires WebView2 runtime (pre-installed on Windows 10/11)

## How It Works

### Architecture Overview

GoDesktop uses a two-stage architecture for maximum efficiency:

#### Stage 1: Build Time (CLI Creation)
1. **Platform-specific runners** are pre-compiled for each OS
2. **Runners are embedded** into the main CLI binary using Go's `embed` directive
3. **Single CLI binary** contains all platform targets

#### Stage 2: App Creation (Runtime)
1. **Instant app generation**: Embedded runner is written to disk
2. **Configuration injection**: App settings stored in JSON (macOS) or passed as flags (Windows)
3. **Platform-specific packaging**:
   - **macOS**: Creates `.app` bundle with proper directory structure, Info.plist, and icon conversion
   - **Windows**: Creates standalone `.exe` with embedded config


### Technical Chain

```
┌─────────────────┐     ┌──────────────────┐    ┌─────────────────┐
│   Your Website  │     │   GoDesktop CLI  │    │  Native Desktop │
│                 │───▶│                  │───▶│      App        │
│ https://...     │     │  Embedded Runner │    │   WebView GUI   │
└─────────────────┘     └──────────────────┘    └─────────────────┘
```

**Why this approach works:**

🎯 **No Runtime Dependencies**: Apps use the system's native web engine
- macOS: WebKit (Safari engine)
- Windows: WebView2 (Chromium Edge engine)


### Size Comparison

| Approach | Bundle Size | Runtime |
|----------|-------------|---------|
| **GoDesktop** | **~3-8MB** | **System WebView** |
| Electron | ~100MB | Bundled Chromium |

The resulting apps are **completely self-contained** and behave like native applications

## License
MIT
