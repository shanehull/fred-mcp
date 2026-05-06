# FRED MCP Server — Implementation Plan

## Overview

Build an MCP (Model Context Protocol) server that wraps the
[go-fred](https://github.com/shanehull/go-fred) FRED API client library,
exposing all 37 FRED API endpoints as MCP tools. The server supports three
deployment modes:

1. **Stdio binary** — Runs locally via `fred-mcp` or
   `go run github.com/shanehull/fred-mcp/cmd/fred-mcp@latest`, connecting
   over stdio (no auth needed — API key via env var).
2. **Docker container** — HTTP/SSE/Streamable HTTP server with dynamic OAuth
   auth, published to GitHub Packages.
3. **Server binary** — Standalone HTTP server (`fred-mcp serve`) with the
   same dynamic auth.

---

## Goals

- Single binary (`fred-mcp`) with two modes: stdio (default) and serve
  subcommand.
- Expose all 37 FRED API endpoints from `go-fred` as typed MCP tools.
- Server mode with dynamic OAuth authentication (RFC 9728 / RFC 8414),
  patterned after `obsidian-remote`.
- Stdio mode for local CLI / Claude Desktop / Codex usage.
- Docker image published to `ghcr.io/shanehull/fred-mcp`.
- Cross-compiled binaries (linux/amd64, linux/arm64, darwin/amd64,
  darwin/arm64) published on GitHub Releases.
- Release automation via release-please (like `go-fred`) with binary
  uploads (like `shed`).

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                      MCP Clients                                  │
│  (Claude Desktop, Codex, Cursor, Gemini CLI, Goose, etc.)        │
└──────────────┬──────────────────────┬────────────────────────────┘
               │                      │
          stdio (local)         HTTP / SSE (remote)
               │                      │
               ▼                      ▼
┌──────────────────────────────────────────────────────────────────┐
│                    cmd/fred-mcp/main.go                            │
│                                                                    │
│  ┌─────────────────────┐    ┌─────────────────────────────────┐  │
│  │ fred-mcp (default)  │    │ fred-mcp serve (:4000)          │  │
│  │                     │    │                                 │  │
│  │ MCPServer (stdio)   │    │ Middleware: CORS + Logging      │  │
│  │ • FRED_API_KEY env  │    │ Auth Middleware (JWT+tokeninfo) │  │
│  │ • 37 tools          │    │ MCPServer (SSE + Streamable)    │  │
│  │ • serve.Stdio()     │    │                                 │  │
│  └─────────────────────┘    │ OAuth Proxy Routes:             │  │
│                              │ • /.well-known/...             │  │
│                              │ • /register, /authorize, /token│  │
│                              │ • /config                      │  │
│                              │ • /sse, /message, /mcp         │  │
│                              └───────────────┬─────────────────┘  │
│                                              │                    │
│                              ┌───────────────▼─────────────────┐  │
│                              │ fred.Client (go-fred library)   │  │
│                              └─────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
fred-mcp/
├── cmd/
│   └── fred-mcp/
│       └── main.go            # Single entry point (stdio + serve subcommand)
├── internal/
│   ├── config/
│   │   └── config.go          # Environment-based configuration
│   ├── handlers/
│   │   ├── tools.go           # Tool registration (RegisterTools)
│   │   └── http.go            # OAuth proxy / discovery handlers (serve mode)
│   ├── middleware/
│   │   └── auth.go            # Auth middleware (serve mode)
│   └── tools/
│       ├── series.go          # Series info, categories, release, tags (simple tools)
│       ├── observation.go     # Series observations, vintages, as-of (complex tools)
│       ├── category.go        # Category tool implementations
│       ├── release.go         # Release tool implementations
│       ├── search.go          # Search + cross-reference tool implementations
│       ├── source.go          # Source tool implementations
│       ├── tags.go            # Tag tool implementations
│       └── geofred.go         # GeoFRED tool implementations
├── Dockerfile                  # Multi-stage build (go 1.25 → scratch)
├── .github/
│   └── workflows/
│       ├── lint.yaml          # golangci-lint on PRs
│       ├── test.yaml          # go test on PRs + main
│       ├── release.yaml       # release-please + binary build & upload
│       ├── docker.yaml        # Docker build & push to ghcr.io
│       ├── nix-eval.yaml      # Verify flake evaluates on all systems (PRs)
│       └── nix-hashes.yaml    # Auto-update vendorHash when go.sum changes
├── flake.nix                   # Nix flake (package + overlay)
├── nix/
│   ├── package.nix             # buildGoModule derivation
│   └── hashes.json             # Single vendorHash (platform-independent for Go)
├── release-please-config.json
├── .release-please-manifest.json
├── go.mod
├── go.sum
├── AGENTS.md
├── README.md
└── PLAN.md                    # This file
```

---

## Single Binary: `fred-mcp`

One binary at `cmd/fred-mcp/main.go` with two modes:

### Default mode — Stdio

```
fred-mcp
FRED_API_KEY=xxx fred-mcp
go run github.com/shanehull/fred-mcp/cmd/fred-mcp@latest
```

- Reads `FRED_API_KEY` from environment.
- Creates an MCPServer, registers all 37 tools, serves over stdio.
- No auth middleware.

Claude Desktop config:

```json
{
  "mcpServers": {
    "fred": {
      "command": "fred-mcp",
      "env": { "FRED_API_KEY": "your-fred-api-key" }
    }
  }
}
```

### Serve mode — HTTP Server

```
fred-mcp serve
FRED_API_KEY=xxx fred-mcp serve
```

- Reads full configuration from environment (see below).
- Creates an MCPServer with SSE + Streamable HTTP transports.
- If `OAUTH_AUDIENCE` is set → JWT/OAuth auth middleware on MCP routes.
- If `OAUTH_AUDIENCE` is unset → all routes are unauthenticated (dev mode).
- Serves OAuth proxy routes for dynamic client auth.

### Docker mode

```
docker pull ghcr.io/shanehull/fred-mcp:latest
docker run -p 4000:4000 \
  -e FRED_API_KEY=xxx \
  -e OAUTH_AUDIENCE=xxx \
  ghcr.io/shanehull/fred-mcp:latest
```

The Docker image runs the serve mode by default (ENTRYPOINT is already `["fred-mcp", "serve"]`). What the image does not include is the FRED API key — that must be supplied via `-e FRED_API_KEY`. Similarly OAuth vars.

---

## Tool Mapping (go-fred → MCP Tools)

37 MCP tools, one per FRED API endpoint. Each tool:

1. Extracts typed parameters from the MCP request.
2. Calls the corresponding `go-fred` method.
3. Returns the result as JSON text.

### Series Tools (12)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_series_info` | `GetSeriesInfo` | `series_id` |
| `fred_get_series_observations` | `GetSeriesObservations` | `series_id`, observation options |
| `fred_get_series_all_releases` | `GetSeriesAllReleases` | `series_id`, observation options |
| `fred_get_series_first_release` | `GetSeriesFirstRelease` | `series_id`, observation options |
| `fred_get_series_as_of` | `GetSeriesAsOf` | `series_id`, `as_of_date` |
| `fred_get_series_vintage_dates` | `GetSeriesVintageDates` | `series_id`, observation options |
| `fred_get_series_categories` | `GetSeriesCategories` | `series_id` |
| `fred_get_series_release` | `GetSeriesRelease` | `series_id` |
| `fred_get_series_tags` | `GetSeriesTags` | `series_id`, tag options |
| `fred_search_series_tags` | `SearchSeriesTags` | `series_search_text`, tag options |
| `fred_search_series_related_tags` | `SearchSeriesRelatedTags` | `series_search_text`, `tag_names` |
| `fred_get_series_updates` | `GetSeriesUpdates` | update options |

### Category Tools (5)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_category` | `GetCategory` | `category_id` |
| `fred_get_category_children` | `GetCategoryChildren` | `category_id` |
| `fred_get_category_related` | `GetCategoryRelated` | `category_id` |
| `fred_get_category_tags` | `GetCategoryTags` | `category_id`, tag options |
| `fred_get_category_related_tags` | `GetCategoryRelatedTags` | `category_id`, `tag_names` |

### Release Tools (8)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_releases` | `GetReleases` | release list options |
| `fred_get_releases_dates` | `GetReleasesDates` | release date options |
| `fred_get_release` | `GetRelease` | `release_id` |
| `fred_get_release_dates` | `GetReleaseDates` | `release_id`, release date options |
| `fred_get_release_sources` | `GetReleaseSources` | `release_id` |
| `fred_get_release_tags` | `GetReleaseTags` | `release_id`, tag options |
| `fred_get_release_related_tags` | `GetReleaseRelatedTags` | `release_id`, `tag_names` |
| `fred_get_release_tables` | `GetReleaseTables` | `release_id`, table options |

### Source Tools (3)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_sources` | `GetSources` | source options |
| `fred_get_source` | `GetSource` | `source_id` |
| `fred_get_source_releases` | `GetSourceReleases` | `source_id`, release list options |

### Tag Tools (3)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_tags` | `GetTags` | tag options |
| `fred_get_related_tags` | `GetRelatedTags` | `tag_names`, tag options |
| `fred_get_tags_series` | `GetTagsSeries` | `tag_names`, search options |

### Search Tools (3)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_search_series` | `SearchSeries` | `search_text`, search options |
| `fred_get_release_series` | `GetReleaseSeries` | `release_id`, search options |
| `fred_get_category_series` | `GetCategorySeries` | `category_id`, search options |

### GeoFRED Tools (3)

| MCP Tool | go-fred Method | Key Parameters |
|---|---|---|
| `fred_get_series_group` | `GetSeriesGroup` | `series_id` |
| `fred_get_series_data` | `GetSeriesData` | `series_id`, map data options |
| `fred_get_regional_data` | `GetRegionalData` | regional data options |

---

## Configuration

All configuration via environment variables, patterned after `obsidian-remote`:

### Common (both modes)

| Variable | Default | Purpose |
|---|---|---|
| `FRED_API_KEY` | *(required)* | FRED API key for data access |

### Serve mode only

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `4000` | HTTP listen port |
| `OAUTH_ISSUER` | `https://accounts.google.com` | OAuth issuer URL |
| `OAUTH_JWKS_URL` | `https://www.googleapis.com/oauth2/v3/certs` | JWKS endpoint |
| `OAUTH_AUTHORIZE_URL` | `https://accounts.google.com/o/oauth2/v2/auth` | OAuth authorize endpoint |
| `OAUTH_TOKEN_URL` | `https://oauth2.googleapis.com/token` | OAuth token endpoint |
| `PUBLIC_HOST` | *(empty)* | Public host (for discovery URLs, SSE base URL) |
| `OAUTH_AUDIENCE` | *(empty)* | OAuth client ID; unset = no auth |
| `OAUTH_CLIENT_SECRET` | *(empty)* | OAuth client secret (server-side) |
| `OAUTH_ALLOWED_EMAIL` | *(empty)* | Restrict to specific email |

---

## Dynamic Auth (Serve Mode)

Follows the same OAuth proxy pattern as `obsidian-remote`:

1. **RFC 9728 discovery** at `/.well-known/oauth-protected-resource/...`
2. **RFC 8414 discovery** at `/.well-known/oauth-authorization-server`
3. **Dynamic client registration** at `/register` — returns server's `client_id`
4. **Authorize proxy** at `/authorize` — injects `scope` parameter
5. **Token proxy** at `/token` — injects `client_id` + `client_secret`
6. **Dynamic config** at `/config` — returns `{type:"oauth", issuer:"...", clientId:"..."}`
7. **Auth middleware** — JWT validation via JWKS, opaque token fallback via
   Google `tokeninfo`

MCP routes (auth-gated when `OAUTH_AUDIENCE` is set):
- `/sse` — SSE transport
- `/message` — SSE message handler
- `/mcp` — Streamable HTTP transport

When `OAUTH_AUDIENCE` is unset, JWKS is never fetched and the auth middleware
is a no-op pass-through (dev mode).

---

## Docker Image

Multi-stage build (simpler than `obsidian-remote` since no Obsidian desktop needed):

```dockerfile
# Stage 1: Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/fred-mcp ./cmd/fred-mcp

# Stage 2: Runtime
FROM scratch
COPY --from=builder /app/fred-mcp /usr/local/bin/fred-mcp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV PORT=4000
EXPOSE 4000
ENTRYPOINT ["/usr/local/bin/fred-mcp", "serve"]
```

Published to `ghcr.io/shanehull/fred-mcp` via GitHub Actions.

Tags:
- `main` — latest main branch
- `v1.2.3` — semver tag
- `v1.2` — minor version
- `latest` — latest release

---

## Nix Flake

Provides a `flake.nix` (patterned after `opencode`) for installing `fred-mcp`
via Nix. The flake exports:

- **`packages.${system}.default`** — the `fred-mcp` binary built with
  `buildGoModule`.
- **`overlays.default`** — an overlay adding `fred-mcp` to `pkgs`.
- **`devShells.${system}.default`** — a dev shell with Go 1.25 and
  golangci-lint.

### Installation via Nix

```bash
# Direct install
nix profile install github:shanehull/fred-mcp

# As a flake input with overlay
inputs.fred-mcp.url = "github:shanehull/fred-mcp";
environment.systemPackages = [ inputs.fred-mcp.packages.${system}.default ];
```

### flake.nix

```nix
{
  description = "FRED MCP server — MCP tools for Federal Reserve Economic Data";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "aarch64-linux"
        "x86_64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];
      forEachSystem = f: nixpkgs.lib.genAttrs systems
        (system: f nixpkgs.legacyPackages.${system});
      rev = self.shortRev or self.dirtyShortRev or "dirty";
    in
    {
      devShells = forEachSystem (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [ go_1_25 golangci-lint ];
        };
      });

      overlays.default = final: _prev: {
        fred-mcp = final.callPackage ./nix/package.nix { };
      };

      packages = forEachSystem (pkgs: {
        default = pkgs.callPackage ./nix/package.nix { };
      });
    };
}
```

### nix/package.nix

```nix
{ lib, buildGoModule }:

let
  manifest = lib.importJSON ../.release-please-manifest.json;
  hashes = lib.importJSON ./hashes.json;
in
buildGoModule {
  pname = "fred-mcp";
  version = manifest.".";
  src = ../.;
  vendorHash = hashes.vendorHash;
  subPackages = [ "cmd/fred-mcp" ];

  ldflags = [
    "-s" "-w"
    "-X main.version=${manifest."."}"
  ];
  CGO_ENABLED = 0;

  meta = with lib; {
    description = "MCP server for Federal Reserve Economic Data (FRED)";
    homepage = "https://github.com/shanehull/fred-mcp";
    license = licenses.mit;
    mainProgram = "fred-mcp";
    platforms = [
      "aarch64-linux"
      "x86_64-linux"
      "aarch64-darwin"
      "x86_64-darwin"
    ];
  };
}
```

### Keeping vendorHash in sync

Since Go `vendorHash` is platform-independent, we store a single hash in
`nix/hashes.json` (simpler than opencode's 4-platform approach for
non-native `node_modules`):

```json
{
  "vendorHash": "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
}
```

A `nix-hashes` workflow (like opencode's `nix-hashes.yml`) triggers on pushes
to `main` when `go.sum` changes:

1. Runs `nix build .` with `lib.fakeHash` as `vendorHash`.
2. The build fails with a hash mismatch, printing the correct hash.
3. Extracts the correct hash from the error output:
   ```bash
   nix build . 2>&1 | grep -oP 'got:\s*\Ksha256-[A-Za-z0-9+/=]+' | tail -n1
   ```
4. Writes it to `nix/hashes.json` and commits the update.

Since `go.sum` only changes when dependencies change, this workflow runs
infrequently — typically only when `go-fred` or `mcp-go` are upgraded.

The `nix/package.nix` reads from `nix/hashes.json` at evaluation time, so the
hash stays in sync automatically.

### Zero manual maintenance

| What | How it stays current |
|---|---|
| **version** | Read from `.release-please-manifest.json` at eval time (release-please auto-updates it). Injected into the binary via `-ldflags="-X main.version=..."` |
| **vendorHash** | Auto-updated by `nix-hashes.yaml` workflow when `go.sum` changes |
| **go.sum** | Managed by `go mod tidy` / Dependabot |
| **Code changes** | No effect on either — `vendorHash` only covers deps, not project source |

The version string (`v0.2.0`, etc.) is surfaced to users in two places:
- **MCP `initialize` response** — `server.NewMCPServer("FRED MCP", version)` reports it to every client. Claude Desktop, Codex, etc. display it in their server list.
- **`--version` flag** — `fred-mcp --version` prints the version and exits.

### nix-eval workflow

A lightweight `nix-eval` workflow (like opencode's) runs on PRs, verifying the
flake evaluates on all 4 systems without actually building:

```bash
nix flake show --all-systems
for system in x86_64-linux aarch64-linux x86_64-darwin aarch64-darwin; do
  nix eval ".#packages.$system.default.drvPath" --raw
done
```

---

## Release Process

Combines `go-fred`'s release-please approach with `shed`'s binary publishing:

### release-please-config.json

```json
{
  "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
  "release-type": "go",
  "packages": {
    ".": {
      "changelog-path": "CHANGELOG.md",
      "bump-minor-pre-major": true,
      "bump-patch-for-minor-pre-major": true
    }
  }
}
```

### Release Workflow (`release.yaml`)

Triggers on push to `main`:

1. **release-please-action** — analyzes conventional commits, creates release
   PRs, cuts GitHub Releases with semver tags.
2. **If release created** — checkout → setup Go → cross-compile with xgo
   (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64) → generate SHA-256
   checksums → upload binaries (`fred-mcp-linux-amd64`, etc.) to the GitHub Release.

### Docker Workflow (`docker.yaml`)

Triggers on push to `main` and on tag creation (`v*`):

1. Build multi-platform Docker image (linux/amd64, linux/arm64).
2. Push to `ghcr.io/shanehull/fred-mcp`.
3. Skip if commit message contains `chore(main): release` (avoids duplicate build).

---

## Implementation Phases

### Phase 1: Project Scaffold

- [ ] Initialize Go module (`github.com/shanehull/fred-mcp`)
- [ ] Add dependencies: `go-fred`, `mcp-go`, `keyfunc`, `golang-jwt`
- [ ] Create directory structure
- [ ] Implement `internal/config/config.go`
- [ ] Write `AGENTS.md` (conventions, commands, gotchas — patterned after go-fred's)

### Phase 2: MCP Tools

- [ ] Implement all 37 tool handlers in `internal/tools/`
- [ ] Implement `internal/handlers/tools.go` — `RegisterTools()`
- [ ] Each tool: extract params → call go-fred → return JSON

### Phase 3: Stdio Mode

- [ ] Implement `cmd/fred-mcp/main.go` — default (stdio) path
- [ ] Wire up MCPServer with stdio transport
- [ ] Test with Claude Desktop config

### Phase 4: Serve Mode

- [ ] Implement `fred-mcp serve` subcommand in `cmd/fred-mcp/main.go`
- [ ] Implement `internal/handlers/http.go` (OAuth proxy routes)
- [ ] Implement `internal/middleware/auth.go`
- [ ] Wire up SSE + Streamable HTTP transports
- [ ] Test with `OAUTH_AUDIENCE` set and unset

### Phase 5: Docker

- [ ] Write `Dockerfile` (multi-stage, scratch base)
- [ ] Write `.github/workflows/docker.yaml`
- [ ] Test local build: `docker build -t fred-mcp . && docker run --rm fred-mcp serve`
- [ ] Push to ghcr.io

### Phase 6: Nix

- [ ] Write `flake.nix`
- [ ] Write `nix/package.nix` with `buildGoModule`
- [ ] Write `nix/hashes.json` (single vendorHash — platform-independent for Go)
- [ ] Write `.github/workflows/nix-eval.yaml` (flake eval on PRs)
- [ ] Write `.github/workflows/nix-hashes.yaml` (auto-update hash when go.sum changes)
- [ ] Test: `nix build .` and `nix profile install .`

### Phase 7: Release Pipeline

- [ ] Configure release-please
- [ ] Write `.github/workflows/release.yaml` with xgo cross-compilation
- [ ] Write `.github/workflows/test.yaml`
- [ ] Write `.github/workflows/lint.yaml`

### Phase 8: Documentation

- [ ] Write `README.md` with:
  - Quick start (all 3 modes: stdio, serve, Docker)
  - Tool reference (all 37 tools with examples)
  - Configuration reference
  - Claude Desktop setup
  - Docker setup
  - OAuth setup (for serve mode)

---

## Dependencies

| Module | Purpose |
|---|---|
| `github.com/shanehull/go-fred` | FRED API client |
| `github.com/mark3labs/mcp-go` | MCP server framework |
| `github.com/MicahParks/keyfunc/v2` | JWKS key fetching (serve mode) |
| `github.com/golang-jwt/jwt/v5` | JWT parsing (serve mode) |

---

## Testing Strategy

Follows `go-fred` conventions: stdlib `testing` only, external test packages
(`_test` suffix), no testify or gomega.

### Layer 1: Tool handler unit tests (always run, no API key)

Every tool handler in `internal/tools/` gets a `*_test.go` file. Each test
starts a local `httptest.Server` that returns known JSON, constructs a
`fred.Client` pointed at it, then calls the handler and verifies the result.

```go
// internal/tools/series_test.go
package tools_test

import "net/http/httptest"

func TestHandleGetSeriesInfo(t *testing.T) {
    mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify request path/params
        if r.URL.Query().Get("series_id") != "GDP" {
            t.Error("expected series_id=GDP")
        }
        w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product"}]}`))
    }))
    defer mock.Close()

    client, _ := fred.New(
        fred.WithAPIKey("test"),
        fred.WithHTTPClient(mock.Client()),
        fred.WithBaseURL(mock.URL),
    )

    req := mcp.CallToolRequest{...}
    result, err := handleGetSeriesInfo(context.Background(), client, req)
    if err != nil { t.Fatal(err) }
    // verify result contains expected JSON
}
```

**Mandated coverage per test file:**

| Test | What it verifies |
|---|---|
| Happy path | Correct params forwarded, valid JSON returned |
| Missing required param | Returns error, not panic |
| FRED API error response | Error propagated correctly |
| Nested options | ObservationOption/TagOption/SearchOption params mapped correctly |
| Empty result | Handles zero results gracefully |
| Complex types | Date parsing, numeric values, nested structures |

Every tool must have at least: happy path + missing param + API error.

### Layer 2: Integration tests (gated by `FRED_API_KEY`)

Patterned after `go-fred`: tests hit the live FRED API, `t.Skip()` if the
env var is absent. At least one tool per category gets an integration test.

```go
func TestHandleGetSeriesObservations_Live(t *testing.T) {
    if os.Getenv("FRED_API_KEY") == "" {
        t.Skip("FRED_API_KEY not set")
    }
    client, _ := fred.New()
    // ...call handler, verify non-empty result
}
```

**Mandated:** one integration test per tool category (7 total: series, category,
release, source, tag, search, geofred).

### Layer 3: Server tests (always run, no API key)

| Test file | Coverage |
|---|---|
| `internal/config/config_test.go` | All env var parsing, defaults, required vars |
| `internal/middleware/auth_test.go` | Valid JWT passes; invalid/expired JWT returns 401; missing Bearer returns 401 with WWW-Authenticate; `OAUTH_AUDIENCE` unset → no-op pass-through |
| `internal/handlers/http_test.go` | OAuth proxy routes return correct 3xx/JSON; discovery endpoints return valid RFC 9728 / 8414 metadata |

### Layer 4: CLI tests (optional)

Verify `--version` outputs version string, `serve` subcommand wires routes.

### What we do NOT test

- **go-fred internals** — already covered by go-fred's own 38 tests.
- **mcp-go internals** — covered by the framework.
- **Network timeouts / retries** — go-fred handles this.

### Running tests

```bash
# Unit tests (no API key needed)
go test ./...

# Integration tests (needs FRED API key)
FRED_API_KEY=xxx go test ./... -run Live

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## CI Workflows

### `lint.yaml`
- Trigger: PRs to `main`
- Runs: `golangci-lint`

### `test.yaml`
- Trigger: PRs to `main`, pushes to `main` (Go files only)
- Jobs:
  - `unit`: `go test -short ./...` (skips live API tests)
  - `integration`: `FRED_API_KEY=${{ secrets.FRED_API_KEY }} go test -run Live ./...` (skips if secret absent)

### `release.yaml`
- Trigger: Push to `main`
- Steps: release-please → (if release) xgo build → upload binaries

### `docker.yaml`
- Trigger: Push to `main`, tag `v*`
- Steps: Docker build → push to ghcr.io
- Skips release-please commits

### `nix-eval.yaml`
- Trigger: PRs to `main`
- Steps: `nix flake show --all-systems` + `nix eval` drvPath for each system
- Verifies the flake evaluates without building

### `nix-hashes.yaml`
- Trigger: Push to `main` when `go.sum` or `nix/package.nix` changes
- Steps: Build with `lib.fakeHash` → extract correct vendorHash from error
  output → commit to `nix/hashes.json`
- Patterned after opencode's `nix-hashes.yml`
