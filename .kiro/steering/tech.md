# Technology Stack

## Core Technologies
- **Go 1.24**: Primary programming language
- **go-rod**: Browser automation library for headless Chrome/Chromium control
- **fsnotify**: File system event monitoring for watch mode
- **Embedded assets**: Frontend assets bundled using Go embed directive

## Dependencies
- `github.com/go-rod/rod v0.116.2` - Browser automation
- `github.com/fsnotify/fsnotify v1.9.0` - File watching
- `github.com/ysmood/gson v0.7.3` - JSON handling

## Frontend Assets
- **Structurizr UI**: Complete web-based diagram renderer
- **Joint.js 3.6.5**: Diagramming library
- **jQuery 3.6.3**: DOM manipulation
- **Bootstrap 3.3.7**: UI framework
- **Backbone.js 1.4.1**: MVC framework
- **Lodash 4.17.21**: Utility library

## Build System

### Common Commands
```bash
# Build for current platform
go build -o bin/strender

# Build all platforms (via script)
./scripts/build.sh

# Install locally
./scripts/install.sh

# Development with watch mode
./strender -watch workspace.json

# Debug mode (non-headless)
./strender -debug workspace.json
```

### Build Configuration
- Cross-compilation for multiple platforms (Darwin/Linux, ARM64/AMD64)
- CGO disabled for static binaries
- Optimized builds with `-ldflags="-s -w"` for size reduction
- Binaries output to `bin/` directory with platform-specific naming

### Environment Variables
- `PORT`: HTTP server port (default: 3000) for serving embedded frontend