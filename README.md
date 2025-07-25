# High-Performance-Event-Driven-Proxy-Server

Production-ready event-driven proxy server written in Go.

## Architecture

```mermaid
flowchart TD
    subgraph User
        A[Client]
    end
    subgraph ProxyServer
        B["Listener HTTP/TCP/UDP"]
        C["Core Engine"]
        D["Config Loader"]
        E["Service Discovery"]
        F["gRPC Admin API"]
        G["Plugin System (Lua/WASM)"]
        H["Metrics & Tracing"]
    end
    subgraph Backends
        I["Backend Services"]
    end
    subgraph External
        J["Consul"]
        K["Prometheus/Grafana"]
        L["OpenTelemetry Collector"]
    end

    A -->|HTTP/TCP/UDP| B
    B --> C
    C -->|middleware| G
    G --> C
    C -->|load balancing| I
    C --> H
    C --> F
    C --> D
    D -->|hot reload| C
    D --> E
    E -->|update backends| C
    E --> J
    H --> K
    H --> L
    F -->|manage| D
    F -->|manage| E
```

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

## Docker
Build and run:
```
docker build -t proxy-server .
docker run -p 8080:8080 -p 9000:9000 -p 9001:9001 -p 9100:9100 -p 9090:9090 proxy-server
```

## Testing
Run all integration tests:
```
go test ./test/ -v
```

## Troubleshooting
- **Port already in use:**
  - Stop any process using the port (e.g. `lsof -i :8080` and `kill <PID>`), or run Docker with different ports.
- **Consul connection refused:**
  - Start Consul locally: `docker run -d --name=consul -p 8500:8500 consul`
  - Or disable discovery in config if not needed.

## License
MIT 