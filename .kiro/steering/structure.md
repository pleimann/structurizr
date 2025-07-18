# Project Structure

## Root Directory
- `main.go` - Entry point with CLI argument parsing and file watching
- `go.mod` / `go.sum` - Go module definition and dependencies
- `big-bank-plc.json` - Example Structurizr workspace file
- `LICENSE` - Apache 2.0 license

## Core Packages
- `renderer/` - Main rendering logic package
  - `renderer.go` - Core renderer implementation with browser automation

## Frontend Assets
- `renderer/frontend/` - Embedded web UI assets (served via Go embed)
  - `index.html` - Main HTML entry point
  - `css/` - Stylesheets for Structurizr UI
  - `js/` - JavaScript libraries and Structurizr implementation

## Build & Deployment
- `bin/` - Compiled binaries (generated)
- `scripts/` - Build and deployment scripts
  - `build.sh` - Cross-platform compilation
  - `install.sh` - Local installation
  - `release.sh` - Release automation
  - `structurizr-pull.sh` - Workspace synchronization

## Code Organization Patterns

### Package Structure
- Single main package for CLI entry point
- `renderer` package encapsulates all rendering logic
- Embedded assets using `//go:embed` directive

### File Naming
- Go files use snake_case when multiple words needed
- Scripts use kebab-case with `.sh` extension
- Binaries follow pattern: `strender-{platform}-{arch}`

### Key Architectural Patterns
- Embedded HTTP server for serving frontend assets
- Browser automation via go-rod for headless rendering
- File system watching for development workflow
- Cross-platform build system with shell scripts