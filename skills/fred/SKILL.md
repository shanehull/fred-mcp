---
name: fred
description: Fetch economic data from the Federal Reserve — GDP, unemployment, inflation, interest rates, and any FRED/ALFRED/GeoFRED series. Use this skill when the user asks about economic indicators, historical trends, or wants to browse Federal Reserve releases by category, even if they don't mention FRED by name.
compatibility: Requires FRED MCP server connection and a free FRED API key.
allowed-tools: fred_*
---

# FRED — Federal Reserve Economic Data

Read-only MCP tools for the Federal Reserve Economic Data API.

## Parameter Formatting

- **IDs as strings** — Category, release, and source IDs are strings (`"1"`, `"53"`).
- **Tag names** — Comma-separated (`"gdp,inflation"`).
- **Dates** — YYYY-MM-DD format.
- **Limit** — Number type for observation tools, string type for search/tag tools.

## Workflows

**Find a series**: `search_series` → pick from results → `get_series_info` → `get_series_observations`.

**Latest data point**: `get_series_observations` with `sort_order="desc"` and `limit="1"`.

**Historical revisions**: `get_series_all_releases` with `realtime_start`/`realtime_end` filters (narrow the window — 1776–9999 exceeds the API limit).

**Category browse**: `get_category_children` starting at `"0"` (root) → drill down → `get_category_series`.

**Vintage data**: `get_series_as_of` with `as_of_date` to see values as-known at a point in history.

## Gotchas

- **Series IDs are case-sensitive** — `gdP` won't match `GDP`. Always search first with `search_series` to confirm the exact ID.
- **Release date range limit** — `realtime_start` to `realtime_end` must be within a ~30-year window. The API rejects `1776-01-01` to `9999-12-31`. Narrow to the release vintages you actually need.
- **Observation count limits** — the API caps observations per call. Use `limit` and `sort_order="desc"` for latest values; paginate with `offset` for long histories.
- **GeoFRED series** use the same endpoints (`get_series_observations`, etc.) — no separate API. Just pass the series ID (e.g., `WIPCPI`).

## Reference IDs

| Type     | ID          | Name                                      |
| -------- | ----------- | ----------------------------------------- |
| Series   | `GDP`       | Gross Domestic Product                    |
| Series   | `UNRATE`    | Unemployment Rate                         |
| Series   | `DGS20`     | 20-Year Treasury Rate                     |
| Series   | `WIPCPI`    | Per Capita Personal Income (GeoFRED)      |
| Category | `0`         | Root                                      |
| Category | `1`         | Production & Business Activity            |
| Release  | `53`        | Gross Domestic Product                    |
| Source   | `1`         | Board of Governors of the Federal Reserve |
| Tag      | `gdp`       | GDP-related                               |
| Tag      | `inflation` | Inflation-related                         |
