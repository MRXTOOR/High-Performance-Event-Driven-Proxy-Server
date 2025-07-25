# High-Performance-Event-Driven-Proxy-Server

Production-ready event-driven proxy server written in Go.

## Features
- HTTP(S)/TCP/UDP proxying
- Horizontal scaling
- Hot reload configuration
- Zero downtime deploy
- gRPC API for management
- Dynamic service discovery (Consul)
- Lua/WASM plugin support
- Prometheus metrics, OpenTelemetry tracing

## Project Structure
```
cmd/proxy-server/      # Entry point
internal/core/         # Proxy engine
internal/config/       # Config loading, hot reload
internal/discovery/    # Service discovery
internal/plugins/      # Lua/WASM plugins
internal/api/          # gRPC API
internal/metrics/      # Metrics and tracing
pkg/                   # Shared libraries
configs/               # Example configs
test/                  # Integration tests
```

## Quick Start
1. Build: `go build -o proxy-server ./cmd/proxy-server`
2. Run: `./proxy-server`
3. Config: see `configs/example.yaml`
4. Metrics: `http://localhost:9100/metrics`
5. gRPC API: `localhost:9090`

## Testing
Run all integration tests:
```
go test ./test/ -v
```

## License
MIT 