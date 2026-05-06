# AGENTS.md — fred-mcp

MCP server wrapping the go-fred FRED API client. 37 tools. Stdio + Serve modes.

## Commands

```bash
go build ./...          # build — must pass
go vet ./...            # vet — must pass
golangci-lint run       # lint — must pass (0 issues)
go test -short ./...    # unit tests — must pass (no API key needed)
go test ./...           # all tests including integration (FRED_API_KEY required for Live tests)
FRED_API_KEY=xxx go test -run Live ./...  # integration tests only
```

Before committing, ensure all four pass:
```bash
go build ./... && go vet ./... && golangci-lint run && go test -short ./...
```

Tool versions are pinned in `mise.toml`.

## Local end-to-end testing

The project includes MCP configs that let you test all 37 tools in a real MCP client:

- **`opencode.json`** — OpenCode config (local stdio). Start with `FRED_API_KEY=xxx opencode` from the project root.
- **`.mcp.json`** — Generic MCP config. Used by Claude Code (`/mcp add`) and other clients that support stdio MCP servers.

Both reference `FRED_API_KEY` from the environment. Set it before launching the client:

```bash
export FRED_API_KEY=your-key
opencode   # loads opencode.json, fred tools auto-register
```

This is useful for verifying tool output, parameter parsing, and error handling without writing automated tests.

## Architecture

```
cmd/fred-mcp/main.go        Single binary entry point (stdio + serve subcommand)
internal/config/config.go   Environment-based configuration
internal/handlers/tools.go  Tool registration (RegisterTools)
internal/handlers/http.go   OAuth proxy / discovery handlers (serve mode)
internal/middleware/auth.go Auth middleware (serve mode)
internal/tools/             Tool handler implementations (37 tools)
```

## Conventions

- Tests in `package_test` (external test package)
- Stdlib `testing` only, no testify or gomega
- Unit tests use `httptest.Server` mock, no live API calls
- Integration tests are gated by `FRED_API_KEY` env var (function names: `TestLive_*`)
- Tools return errors in `CallToolResult`, not as Go errors
- Use `mcp.WithReadOnlyHintAnnotation(true)` and `mcp.WithDestructiveHintAnnotation(false)` for all tools
- Semantic commits for release-please: `feat:` bumps minor (pre-1.0), `fix:` bumps patch. `docs:`, `ci:`, `chore:` don't bump version but appear in changelog.
- Known test series: `DGS20` (20-year Treasury), `UNRATE` (Unemployment Rate), `WIPCPI` (GeoFRED per-capita income), `SMU56000000500000001A` (GeoFRED), release `53` (GDP), category `1` (Production & Business Activity), category `0` (root), source `1` (Federal Reserve), tag `"gdp"`/`"business"`.

## Tool handler pattern

```go
func HandleGetSeriesInfo(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    seriesID, err := req.RequireString("series_id")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    result, err := client.GetSeriesInfo(ctx, seriesID)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
    }
    return tools.MarshalResult(result)
}
```
