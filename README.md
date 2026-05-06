# fred-mcp

<p align="center">
  <img src="https://github.com/modelcontextprotocol/.github/raw/main/profile/assets/light.png" height="80" alt="MCP">
  <br>
  <img src="assets/fred-logo.svg" height="48" alt="FRED">
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shanehull/fred-mcp.svg)](https://pkg.go.dev/github.com/shanehull/fred-mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanehull/fred-mcp)](https://goreportcard.com/report/github.com/shanehull/fred-mcp)
[![CI](https://github.com/shanehull/fred-mcp/actions/workflows/test.yaml/badge.svg)](https://github.com/shanehull/fred-mcp/actions/workflows/test.yaml)
[![Docker](https://img.shields.io/badge/docker-ghcr.io-blue)](https://github.com/shanehull/fred-mcp/pkgs/container/fred-mcp)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

MCP server for [Federal Reserve Economic Data](https://fred.stlouisfed.org/). Wraps all 37 [go-fred](https://github.com/shanehull/go-fred) API endpoints as typed MCP tools. Connects over stdio (local) or HTTP/SSE (remote with OAuth).

## Quick Start

### 1. Get a FRED API key

[Register for free](https://fred.stlouisfed.org/docs/api/api_key.html) — takes 30 seconds.

### 2. Pick a mode

### Stdio

Run locally via the [installed](#install) binary, or with `go run` (no install):

```bash
# installed binary
FRED_API_KEY=your-key fred-mcp

# go run — no install needed
FRED_API_KEY=your-key go run github.com/shanehull/fred-mcp/cmd/fred-mcp@latest
```

Client config for most clients (Claude Code, Codex, Gemini CLI, Goose, VS Code):

```json
{
  "mcpServers": {
    "fred": {
      "command": "fred-mcp",
      "env": { "FRED_API_KEY": "your-key" }
    }
  }
}
```

Or without installing — point your client at `go run`:

```json
{
  "mcpServers": {
    "fred": {
      "command": "go",
      "args": ["run", "github.com/shanehull/fred-mcp/cmd/fred-mcp@latest"],
      "env": { "FRED_API_KEY": "your-key" }
    }
  }
}
```

OpenCode uses its own format:

```jsonc
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "fred": {
      "type": "local",
      "command": ["go", "run", "github.com/shanehull/fred-mcp/cmd/fred-mcp@latest"],
      "enabled": true,
      "environment": {
        "FRED_API_KEY": "{env:FRED_API_KEY}"
      }
    }
  }
}
```

Where to put it: `opencode.json` (OpenCode), `.codex.json` (Codex), `~/.gemini/settings.json` (Gemini CLI), `.vscode/mcp.json` (VS Code). Claude Code users can run `/mcp add` instead.

### Serve

Run as an HTTP server (useful for remote access or Docker):

```bash
FRED_API_KEY=your-key fred-mcp serve
```

Or via Docker:

```bash
docker run -p 4000:4000 \
  -e FRED_API_KEY=your-key \
  ghcr.io/shanehull/fred-mcp:latest
```

Client config for SSE-capable clients (Cursor, etc.):

```json
{
  "mcpServers": {
    "fred": {
      "url": "http://your-host:4000/sse"
    }
  }
}
```

Streamable HTTP clients can use `http://your-host:4000/mcp`.

### 3. Try it

Once connected, ask your agent: _"What's the latest US unemployment rate?"_ — it will call `fred_get_series_observations` with `UNRATE`.

## API Key

Get a free key at [fred.stlouisfed.org/docs/api/api_key.html](https://fred.stlouisfed.org/docs/api/api_key.html). Pass it via the `FRED_API_KEY` environment variable — all modes respect it.


## Serve Mode (HTTP/SSE)

Extended configuration reference for `fred-mcp serve`.

### OAuth (optional)

When `OAUTH_AUDIENCE` is set, MCP routes require a valid Bearer token. The server acts as an OAuth proxy:

| Endpoint | Purpose |
|----------|---------|
| `/.well-known/oauth-protected-resource` | RFC 9728 resource metadata |
| `/.well-known/oauth-authorization-server` | RFC 8414 discovery |
| `/register` | Dynamic client registration |
| `/authorize` | Authorize proxy (injects scope) |
| `/token` | Token proxy (injects client credentials) |
| `/config` | OAuth config endpoint |

Tokens are validated via JWKS (JWT) with opaque token fallback (Google tokeninfo). Restrict access with `OAUTH_ALLOWED_EMAIL`.

### Docker

```bash
docker run -p 4000:4000 \
  -e FRED_API_KEY=your-key \
  ghcr.io/shanehull/fred-mcp:latest
```

The container runs serve mode. Add OAuth env vars to enable authentication. Multi-arch (`linux/amd64`, `linux/arm64`), `FROM scratch` (~8 MB).

| Tag | Description |
|-----|------------|
| `latest` | Latest release |
| `main` | Main branch |
| `v1.2.3` | Specific version |
| `v1.2` | Minor version |

## Tool Reference

All 37 FRED API endpoints, one MCP tool each. Every tool returns JSON.

### Series (12 tools)

| Tool | Description |
|---|---|
| `fred_get_series_info` | Get series metadata (title, units, frequency, tags) |
| `fred_get_series_observations` | Get time-indexed observation data |
| `fred_get_series_all_releases` | Get all vintages (every revision) |
| `fred_get_series_first_release` | Get first-published values (ALFRED) |
| `fred_get_series_as_of` | Get observations as they were on a specific date |
| `fred_get_series_vintage_dates` | List all revision dates |
| `fred_get_series_categories` | Get categories a series belongs to |
| `fred_get_series_release` | Get the release a series belongs to |
| `fred_get_series_tags` | Get tags assigned to a series |
| `fred_search_series_tags` | Search series by tag text |
| `fred_search_series_related_tags` | Find related tags for a series search |
| `fred_get_series_updates` | Get recently updated series |

Observation tools accept optional parameters: `observation_start`, `observation_end`, `realtime_start`, `realtime_end`, `units`, `frequency`, `aggregation_method`, `output_type`, `vintage_dates`, `sort_order`, `limit`, `offset`.

### Categories (5 tools)

| Tool | Description |
|---|---|
| `fred_get_category` | Get a category by ID |
| `fred_get_category_children` | Get child categories (pass 0 for root) |
| `fred_get_category_related` | Get related categories |
| `fred_get_category_tags` | Get tags for a category |
| `fred_get_category_related_tags` | Get related tags for a category |

### Releases (8 tools)

| Tool | Description |
|---|---|
| `fred_get_releases` | List all data releases |
| `fred_get_releases_dates` | List all release dates |
| `fred_get_release` | Get a release by ID |
| `fred_get_release_dates` | Get dates for a release |
| `fred_get_release_sources` | Get sources for a release |
| `fred_get_release_tags` | Get tags for a release |
| `fred_get_release_related_tags` | Get related tags for a release |
| `fred_get_release_tables` | Get tables for a release |

### Sources (3 tools)

| Tool | Description |
|---|---|
| `fred_get_sources` | List all data sources |
| `fred_get_source` | Get a source by ID |
| `fred_get_source_releases` | Get releases for a source |

### Tags (3 tools)

| Tool | Description |
|---|---|
| `fred_get_tags` | Browse/search all tags |
| `fred_get_related_tags` | Find tags related to given tags |
| `fred_get_tags_series` | Get series matching given tags |

### Search (3 tools)

| Tool | Description |
|---|---|
| `fred_search_series` | Search series by keyword |
| `fred_get_release_series` | Get series in a release |
| `fred_get_category_series` | Get series in a category |

Search tools accept optional parameters: `search_type`, `order_by`, `sort_order`, `filter_variable`, `filter_value`, `limit`, `tag_names`, `exclude_tag_names`.

### GeoFRED (3 tools)

| Tool | Description |
|---|---|
| `fred_get_series_group` | Get GeoFRED series group metadata |
| `fred_get_series_data` | Get map data for a series |
| `fred_get_regional_data` | Get regional map data by series group |

## Configuration

All settings via environment variables.

### Common

| Variable | Default | Purpose |
|---|---|---|
| `FRED_API_KEY` | *(required)* | FRED API key |

### Serve mode

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `4000` | HTTP listen port |
| `OAUTH_ISSUER` | Google | OAuth issuer URL |
| `OAUTH_JWKS_URL` | Google certs | JWKS endpoint |
| `OAUTH_AUTHORIZE_URL` | Google auth | OAuth authorize endpoint |
| `OAUTH_TOKEN_URL` | Google token | OAuth token endpoint |
| `PUBLIC_HOST` | *(auto)* | Public host for discovery URLs |
| `OAUTH_AUDIENCE` | *(empty)* | OAuth client ID; unset = no auth |
| `OAUTH_CLIENT_SECRET` | *(empty)* | OAuth client secret (server-side) |
| `OAUTH_ALLOWED_EMAIL` | *(empty)* | Restrict access to a specific email |

## Install

```bash
# Go install
go install github.com/shanehull/fred-mcp/cmd/fred-mcp@latest

# Nix
nix profile install github:shanehull/fred-mcp

# Docker
docker pull ghcr.io/shanehull/fred-mcp:latest

# Prebuilt binary (from GitHub Releases)
curl -fsSL https://github.com/shanehull/fred-mcp/releases/latest/download/fred-mcp-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o /usr/local/bin/fred-mcp
chmod +x /usr/local/bin/fred-mcp
```
