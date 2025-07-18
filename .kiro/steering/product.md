# Product Overview

## Structurizr Renderer

A command-line tool that renders Structurizr workspace diagrams to PNG images using a headless browser. The tool loads Structurizr workspace JSON files and exports all defined views as individual PNG files.

### Key Features
- Renders Structurizr architecture diagrams to PNG format
- Supports watch mode for automatic re-rendering on file changes
- Headless browser rendering using Chrome/Chromium via go-rod
- Embedded web UI serving Structurizr frontend assets
- Cross-platform binary builds (macOS, Linux, ARM64, AMD64)

### Use Cases
- Generate documentation images from Structurizr workspaces
- Automated diagram generation in CI/CD pipelines
- Local development workflow with auto-refresh capabilities
- Export architecture diagrams for presentations and documentation