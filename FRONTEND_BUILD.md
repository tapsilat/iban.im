# Frontend Build and Embedding Process

## Overview
The frontend is built using Vue 3 + Vite + Tailwind CSS and embedded into the Go binary as static files.

## Frontend Stack
- **Framework**: Vue 3
- **Build Tool**: Vite
- **Styling**: Tailwind CSS v4
- **Router**: Vue Router
- **State Management**: Vuex

## Build Process

**Important**: The `static/dist/` directory is not tracked in git. You must build the frontend and copy it to `static/dist/` before building the Go binary.

### 1. Install Dependencies
```bash
cd web
npm install
```

### 2. Build the Frontend
```bash
cd web
npm run build
```

This creates a `dist` directory with the compiled frontend assets.

### 3. Copy to Static Package
```bash
cp -r web/dist static/
```

The `static/dist` directory contains:
- `index.html` - The main HTML file
- `assets/` - JavaScript and CSS bundles

**Note**: This directory is ignored by git (`.gitignore`) as it contains build artifacts.

### 4. Build the Go Binary
```bash
go build
```

The Go `embed` package automatically includes files from `static/dist/` into the binary.

## How It Works

### Embedding
The `static/embed.go` file uses Go's `embed` directive to include the frontend files:

```go
//go:embed dist
var embedFS embed.FS
```

### Serving
In `main.go`, the embedded files are served:

1. **Static Assets**: `/assets/*` serves JS and CSS files from the embedded filesystem
2. **SPA Routing**: All other routes serve `index.html` to enable client-side routing

### Key Features
- **Zero External Dependencies**: Frontend is compiled into the Go binary
- **Single Binary Deployment**: No need to deploy frontend files separately
- **Client-Side Routing**: Vue Router handles navigation without server round-trips
- **API Integration**: GraphQL API available at `/graph` endpoint

## Development Workflow

### Frontend Development
```bash
cd web
npm run dev
```
This starts a development server on http://localhost:4881 with hot-reload.

### Full Stack Development
1. Start the backend: `go run main.go`
2. Start the frontend dev server: `cd web && npm run dev`
3. Frontend proxies API requests to the backend

## Production Build

### Using Makefile (Recommended)
```bash
make build
```

This single command will:
1. Install frontend dependencies
2. Build the frontend
3. Copy files to `static/dist/`
4. Build the Go binary

### Manual Build
```bash
# Build frontend
cd web
npm run build

# Copy to static package
cp -r dist ../static/

# Build Go binary
cd ..
go build

# Run
./iban.im
```

## Notes
- The frontend is responsive and mobile-friendly
- Tailwind CSS v4 requires `@tailwindcss/postcss` plugin
- All frontend routes are handled client-side except API endpoints
