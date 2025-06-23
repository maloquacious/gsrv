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

See [cmd/example/main.go](cmd/example/main.go) for a complete working example.

To run the example:
```bash
go run cmd/example/main.go
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

## Production Readiness TODO List

### Security & Authentication
- [ ] Rate limiting for shutdown endpoint to prevent abuse
- [ ] More secure shutdown key handling (env vars, key rotation)
- [ ] Input validation and sanitization for all endpoints
- [ ] Add request ID/tracing support for better debugging
- [ ] Security headers middleware (CORS, CSP, etc.)

### Observability & Monitoring  
- [ ] Structured logging with configurable log levels
- [ ] Metrics collection (Prometheus format) - request counts, response times, uptime
- [ ] Custom health check functions (database, external services)
- [ ] Distributed tracing support (OpenTelemetry)
- [ ] Error tracking and alerting integration

### Configuration & Deployment
- [ ] Environment-based configuration (12-factor app)
- [ ] Configuration validation on startup
- [ ] Docker support and multi-stage builds
- [ ] Helm charts for Kubernetes deployment
- [ ] Health check endpoints for container orchestration

### Testing & Quality
- [ ] Comprehensive test coverage (>90%) 
- [ ] Integration tests for graceful shutdown scenarios
- [ ] Benchmarking tests for performance validation
- [ ] Fuzzing tests for security validation
- [ ] CI/CD pipeline with automated testing

### Documentation & Examples
- [ ] Comprehensive API documentation (godoc)
- [ ] Production deployment examples
- [ ] Best practices guide
- [ ] Troubleshooting guide
- [ ] Performance tuning guide

### Error Handling & Resilience
- [ ] Circuit breaker pattern for dependencies
- [ ] Retry logic with exponential backoff
- [ ] Timeout configuration for all operations
- [ ] Graceful degradation strategies
- [ ] Panic recovery middleware

### Performance & Scalability
- [ ] Connection pooling optimizations
- [ ] Memory usage profiling and optimization
- [ ] CPU profiling and bottleneck identification
- [ ] Load testing and capacity planning
- [ ] Horizontal scaling support

### Operations & Maintenance
- [ ] Automated dependency updates
- [ ] Security vulnerability scanning
- [ ] Performance regression testing
- [ ] Rollback strategies
- [ ] Maintenance mode support
