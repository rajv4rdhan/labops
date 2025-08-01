# BunnyPage - Simple Html page with Go Web Server

A minimal Go web server that serves a beautiful animated bunny page. This project demonstrates how to create a lightweight Docker image using Go's static compilation capabilities.

## 🐰 What is BunnyPage?

BunnyPage is a simple HTTP file server written in Go that serves a charming animated webpage featuring a peaceful meadow scene with a cute bunny. It's designed to be deployed as a minimal Docker container using a `scratch` base image.

## 🛠️ Building the Application

### Prerequisites

- Go 1.19 or later
- Docker (for containerization)

### Building the Go Binary

To build a statically linked binary suitable for a `scratch` Docker image:

#### For Windows (PowerShell)

```powershell
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .
```

#### For macOS/Linux

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .
```

### Build Flags Explained

- `CGO_ENABLED=0`: Disables CGO to ensure static linking
- `GOOS=linux`: Target Linux operating system
- `GOARCH=amd64`: Target AMD64 architecture
- `-a`: Force rebuilding of packages
- `-installsuffix cgo`: Use different install suffix to avoid conflicts
- `-ldflags '-extldflags "-static"'`: Link statically
- `-o server`: Output binary name

## 🐳 Docker Deployment

### Building the Docker Image

```bash
docker build -t bunnypage .
```

### Running the Container

```bash
docker run -p 8080:8080 bunnypage
```

### Access the Application

Open your browser and navigate to `http://localhost:8080` to see the bunny page.

## 📁 Project Structure

```text
BunnyPage/
├── main.go          # Go web server source code
├── index.html       # Animated bunny webpage
├── Dockerfile       # Multi-stage Docker build
├── Readme.md        # This file
└── server/          # Directory for built binary (created during build)
```

## 🔧 How It Works

1. **Go Server**: The `main.go` file creates a simple HTTP file server using Go's built-in `net/http` package
2. **Static Files**: Serves the `index.html` file and any other static assets
3. **Docker**: Uses a `scratch` base image for minimal container size
4. **Port**: Listens on port 8080

## 🎨 Features

- **Minimal Size**: Uses `scratch` Docker image for tiny container footprint
- **Static Binary**: No external dependencies required
- **Animated UI**: Beautiful CSS animations and gradients
- **Responsive**: Works on different screen sizes
- **Fast**: Lightweight Go server with minimal overhead

## 📈 Performance Benefits

- **Container Size**: Extremely small Docker image (~10MB total)
- **Memory Usage**: Minimal RAM footprint
- **Startup Time**: Nearly instant container startup
- **Security**: Scratch image reduces attack surface

## 🔧 Development

### Local Development

To run locally for development:

```bash
go run main.go
```

Then visit `http://localhost:8080`

### Testing the Build

Before building the Docker image, test that your binary works:

```bash
# Build the binary
$env:CGO_ENABLED="0"; $env:GOOS="linux"; go build -o server .

# Test it works (on Linux/WSL)
./server
```
