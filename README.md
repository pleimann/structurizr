# Structurizr Renderer GitHub Action

A GitHub Action that renders Structurizr workspace diagrams to PNG images using the structurizr-renderer tool.

## Usage

### Basic Usage

```yaml
name: Render Architecture Diagrams

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  render-diagrams:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Render Structurizr diagrams
        uses: pleimann/structurizr-renderer@v1
        with:
          workspace: 'workspace.json'
          output-dir: 'diagrams'
```

### Advanced Usage

```yaml
name: Generate and Deploy Diagrams

on:
  push:
    branches: [ main ]

jobs:
  render-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Render diagrams
        id: render
        uses: pleimann/structurizr-renderer@v1
        with:
          workspace: 'architecture/workspace.json'
          output-dir: 'docs/images'
          version: 'v1.2.0'
      
      - name: Check render results
        run: |
          echo "Diagrams saved to: ${{ steps.render.outputs.diagrams-path }}"
          echo "Number of diagrams: ${{ steps.render.outputs.diagram-count }}"
      
      - name: Commit generated diagrams
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add docs/images/*.png
          git diff --staged --quiet || git commit -m "Update architecture diagrams"
          git push
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `workspace` | Path to the Structurizr workspace JSON file | Yes | `workspace.json` |
| `output-dir` | Directory to store rendered PNG files | No | `./` |
| `version` | Version of structurizr-renderer to use | No | `latest` |

## Outputs

| Output | Description |
|--------|-------------|
| `diagrams-path` | Path where the rendered diagrams were saved |
| `diagram-count` | Number of diagrams that were rendered |

## Features

- **Cross-platform**: Works on Ubuntu, macOS, and Windows runners
- **Automatic binary download**: Downloads the appropriate binary for the runner's platform
- **Artifact upload**: Automatically uploads generated PNG files as GitHub artifacts
- **Error handling**: Validates workspace file exists and provides clear error messages
- **Flexible versioning**: Use latest release or specify a specific version

## Supported Platforms

- Linux (x64, ARM64)
- macOS (x64, ARM64)

## Example Workflows

### Documentation Generation

```yaml
name: Update Documentation

on:
  push:
    paths:
      - 'architecture/**'
      - 'workspace.json'

jobs:
  update-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Generate diagrams
        uses: pleimann/structurizr-renderer@v1
        with:
          workspace: 'workspace.json'
          output-dir: 'docs/diagrams'
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./docs
```

### Multi-workspace Rendering

```yaml
name: Render Multiple Workspaces

on:
  workflow_dispatch:

jobs:
  render-all:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        workspace:
          - { file: 'system-context.json', dir: 'context' }
          - { file: 'containers.json', dir: 'containers' }
          - { file: 'components.json', dir: 'components' }
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Render ${{ matrix.workspace.file }}
        uses: pleimann/structurizr-renderer@v1
        with:
          workspace: ${{ matrix.workspace.file }}
          output-dir: diagrams/${{ matrix.workspace.dir }}
```

## Requirements

- The workspace JSON file must be valid Structurizr format
- The runner must have internet access to download the binary
- For private repositories, ensure the action has appropriate permissions

## Troubleshooting

### Common Issues

1. **Workspace file not found**: Ensure the workspace path is correct relative to the repository root
2. **Binary download fails**: Check if the specified version exists in the releases
3. **No diagrams generated**: Verify the workspace file contains valid view definitions

### Debug Mode

Add the following step to enable debug output:

```yaml
- name: Enable debug logging
  run: echo "ACTIONS_STEP_DEBUG=true" >> $GITHUB_ENV
```