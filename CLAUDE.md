# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoProxy is a geo-aware proxy gateway written in Go. It stores upstream nodes you control (manual nodes and subscription nodes), validates their reachability and exit metadata (exit IP + geolocation + latency), and selects an available node for each HTTP or SOCKS5 client request. Routing is constrained by region and stabilized with a short-lived session key encoded in the proxy username (username DSL). It exposes one HTTP proxy port, one SOCKS5 proxy port, and an authenticated WebUI.

There is no active public proxy pool or automatic public source collector in this gateway model.

## Build & Run

```bash
# Run directly (requires Go 1.25, CGO enabled for sqlite3)
go run .

# Build and run
go build -o proxygo .
./proxygo

# Docker
docker compose up -d --build
```

CGO is required (`CGO_ENABLED=1`) because of the `github.com/mattn/go-sqlite3` dependency. On Windows a gcc toolchain (e.g. mingw-w64) must be on `PATH`.

## Testing

Go unit tests exist and run with `go test ./...` (CGO required for the `storage`, `custom`, and `webui` packages because they open a real SQLite database):

```bash
go test ./...
go test -race ./affinity ./selector ./auth
```

Key covered paths: username DSL parsing (`auth`), session affinity + TTL GC + race (`affinity`), region selection without cross-region fallback (`selector`), Resolve sticky/rebind decisions (`selector`), storage schema migration + region queries + manual-node CRUD (`storage`), WebUI authentication and no-business-data-leak on unauthenticated requests (`webui`).

Shell/script-based end-to-end tests live under `test/` and run against a live instance with real upstream nodes.

## Architecture

The system is a single binary with several cooperating goroutines. Module `go.mod` name is `goproxy`.

### Module Dependency Flow

```
main.go (orchestrator)
  ├── config/     — Global config (env vars + config.json), thread-safe singleton
  ├── storage/    — SQLite persistence layer (proxies + subscriptions tables)
  ├── auth/       — Username DSL parser + constant-time credential verification
  ├── affinity/   — In-memory session→node binding store with TTL GC
  ├── selector/   — Region + strategy node selection and session Resolve()
  ├── validator/  — Concurrent proxy validation (connectivity + exit IP + geo + latency)
  ├── checker/    — Background health checker (batch-based, skips S-grade when healthy)
  ├── custom/     — Subscription + manual node manager (fetch, parse, validate, refresh)
  │   ├── parser.go   — Clash YAML / plain / base64 / single-link subscription parser
  │   ├── singbox.go  — sing-box process manager (config generation, start/stop/reload)
  │   └── manager.go  — Subscription refresh loop + probe-wake loop for disabled proxies
  ├── proxy/      — Outward-facing proxy servers
  │   ├── server.go        — HTTP proxy (implements http.Handler)
  │   └── socks5_server.go — SOCKS5 proxy (raw TCP, manual protocol implementation)
  ├── webui/      — Dashboard server (embedded HTML in html.go/dashboard.go, API handlers)
  └── logger/     — In-memory log collector for WebUI display
```

### Key Design Patterns

- **Username DSL**: The proxy username encodes routing. Syntax `<base>[-region-<cc>][-session-<id>]`. `region` is a two-letter code normalized to lowercase; `session` is `[A-Za-z0-9_-]` up to 64 chars; order is fixed (region before session). Malformed usernames are rejected with an explicit error, never a silent fallback. Parsed in `auth/dsl.go`.
- **Credential verification**: `auth/credentials.go` uses `subtle.ConstantTimeCompare`. HTTP compares a SHA-256 hash; SOCKS5 compares the plaintext password (protocol constraint). Both share the same base-username check.
- **Region routing without fallback**: `selector.Pick` returns nodes for the requested region only; if none exist it returns `ErrNoNode` (wrapped with the region) and never silently crosses into another region. `DEFAULT_REGION` applies only when the username has no `-region-` suffix.
- **Session affinity**: `affinity.Store` maps a session key to a node binding (`sync.RWMutex`-guarded map). `selector.Resolve` reuses a live, region-matching binding (refreshing LastActive), otherwise picks a new node in the required region and rebinds. A background GC goroutine scans every 1 minute and removes bindings older than `SESSION_TTL_MINUTES` (default 10).
- **Geo-filter**: Whitelist (`AllowedCountries`) takes priority — if non-empty, only those countries pass. Otherwise blacklist (`BlockedCountries`) rejects listed countries. On startup, nodes violating the filter are disabled (not deleted). Applied during validation and startup.
- **Auto-retry on proxy failure**: HTTP and SOCKS5 servers retry with different upstream nodes on failure, excluding already-tried nodes and rebinding sticky sessions as needed. Failed nodes are disabled (soft) rather than deleted, so a probe loop can recover them.
- **SOCKS5 service upstream dialing** supports both SOCKS5 and HTTP upstreams via `dialViaProxy`.
- **Subscription direct fetch**: Subscription URLs are fetched directly (no upstream-pool bounce). Subscriptions with no usable node for 7 days are paused with a warning, not deleted.
- **Encrypted node conversion**: Encrypted subscription/manual single-link nodes are merged into a shared sing-box config, reloaded, and exposed as local SOCKS5 upstreams stored in the DB.

### Background Goroutines (started in main.go)

1. **Session TTL GC** — scans every 1 minute, removes expired session bindings.
2. **Health checker** — every `HEALTH_CHECK_INTERVAL` min, validates a batch of nodes.
3. **Subscription manager** — refresh loop (per-subscription interval) + probe-wake loop for disabled subscription nodes.

### Ports

| Port | Service | Config key |
|------|---------|------------|
| 7800 | WebUI dashboard | `WEBUI_PORT` |
| 7801 | SOCKS5 proxy | `SOCKS5_PORT` |
| 7802 | HTTP proxy | `HTTP_PORT` |

### Configuration

- Environment variables: `WEBUI_PASSWORD`, `PROXY_AUTH_ENABLED`, `PROXY_AUTH_USERNAME`, `PROXY_AUTH_PASSWORD`, `SESSION_TTL_MINUTES`, `DEFAULT_REGION`, `ALLOWED_COUNTRIES`, `BLOCKED_COUNTRIES`, `HEALTH_CHECK_INTERVAL`, `MAX_RETRY`, `SINGBOX_PATH`, `HTTP_PORT`, `SOCKS5_PORT`, `WEBUI_PORT`, `DATA_DIR`
- Persistent config: `config.json` (or `$DATA_DIR/config.json`) — editable via WebUI, overrides env defaults after first save.
- Config is loaded once at startup via `config.Load()`, updated in-memory via `config.Save()`. Thread-safe via `sync.RWMutex`.

### Storage

SQLite with `MaxOpenConns(1)` (single-writer). Two tables: `proxies` (with region/region_source/note metadata and quality grades S/A/B/C based on latency) and `subscriptions` (subscription URL/file metadata, refresh interval, status, last successful refresh). The legacy `source_status` table is dropped on startup. Schema auto-migrates on startup, including region/note/region_source columns and legacy `free`/`custom` source values remapped to `manual`/`subscription`.

### WebUI

The entire frontend is embedded as Go string literals in `webui/html.go` (login) and `webui/dashboard.go` (dashboard). The server (`webui/server.go`) serves HTML and API endpoints. Auth is password-only admin: unauthenticated users see only the neutral login screen, and every business API requires an authenticated session (returns 401 with no business data otherwise). There is no guest/read-only role and no visitor subscription-contribution endpoint.

## Code Conventions

- All log messages use `[module]` prefix: `[main]`, `[health]`, `[custom]`, `[socks5]`, `[proxy]`, `[tunnel]`, `[storage]`, `[config]`, `[webui]`
- Comments and log messages are in Chinese
- Quality grades: S (≤500ms), A (501-1000ms), B (1001-2000ms), C (>2000ms)
- `storage.Proxy` is the shared data type across all modules
- No silent fallbacks: malformed input and missing nodes surface as explicit errors
