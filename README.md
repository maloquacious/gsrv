# gsrv

GSRV implements a graceful shutdown wrapper around Go's `http.Server` that provides:

- Graceful shutdown on OS signals (SIGINT, SIGQUIT, SIGTERM)
- Health check endpoint
- Remote shutdown endpoint with authentication key
- Configurable timeouts and options

## Assumptions

This package assumes you're running behind a reverse proxy (e.g., Nginx) that handles SSL/TLS termination. The server only supports HTTP connections.

## Installation

```bash
go get github.com/maloquacious/gsrv
```

## Usage

### Basic Example

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/maloquacious/gsrv"
)

func main() {
    // Create server with default settings
    server, err := gsrv.New(
        gsrv.WithHost("localhost"),
        gsrv.WithPort("8080"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Set up routes
    mux := http.NewServeMux()
    mux.Handle("GET /health", server.HealthHandler())
    mux.Handle("POST /shutdown/{key}", server.ShutdownHandler())
    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    server.Handler = mux
    
    // Start server - blocks until shutdown signal received
    if err := server.ListenAndServe(); err != nil {
        log.Printf("Server error: %v", err)
    }
}
```

### Configuration Options

```go
server, err := gsrv.New(
    gsrv.WithHost("0.0.0.0"),
    gsrv.WithPort("3000"),
    gsrv.WithShutdownKey("my-secret-key"),
    gsrv.WithContext(ctx),
)
```

### Built-in Handlers

#### Health Check
`GET /health` returns server uptime:
```json
{"uptime": "1h23m45s"}
```

#### Remote Shutdown
`POST /shutdown/{key}` triggers graceful shutdown when the correct key is provided:
```json
{"status": "server shutting down"}
```

The shutdown key is automatically generated if not provided via `WithShutdownKey()`. Retrieve it with `server.ShutdownKey()`.

## Graceful Shutdown

The server automatically handles graceful shutdown:
1. Listens for OS signals (SIGINT, SIGQUIT, SIGTERM)
2. Stops accepting new connections
3. Cancels idle connections
4. Waits up to 10 seconds for active requests to complete
5. Forces shutdown if timeout exceeded
