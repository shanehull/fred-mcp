---
name: fred
description: Access Federal Reserve Economic Data (FRED) via MCP — 37 tools covering series, observations, vintage data, categories, releases, sources, tags, and GeoFRED maps.
compatibility: Requires FRED MCP server connection and a free FRED API key.
allowed-tools: fred_*
---

# FRED — Federal Reserve Economic Data

Query economic time series, search metadata, explore categories, and pull regional map data — all through FRED's API.

## Tools

All tools are read-only. Every tool returns JSON.

### Series & Observations

Fetch time-series data and metadata. Observation tools accept optional parameters for date ranges (`observation_start`, `observation_end`), units (`lin`, `chg`, `pch`), frequency (`m`, `q`, `a`), sort order, and limits.

- `get_series_info` — Series metadata: title, units, frequency, seasonal adjustment, notes.
- `get_series_observations` — Time-indexed data points with dates and values.
- `get_series_updates` — Recently updated series.

### Vintage Data (ALFRED)

FRED preserves every revision. Each data point has `realtime_start` (first reported) and `realtime_end` (last before revision). Use `realtime_start` / `realtime_end` filters to narrow the window.

- `get_series_all_releases` — Every revision, every date.
- `get_series_first_release` — Earliest published value for each date.
- `get_series_as_of` — Data as known on a specific date (provide `as_of_date` in YYYY-MM-DD format).
- `get_series_vintage_dates` — All dates this series was revised.

### Categories

FRED organizes data into categories (topics). Root is `0`, major categories include `1` (Production & Business Activity), etc.

- `get_category` — Get a category by ID.
- `get_category_children` — Child categories (pass `0` for root).
- `get_category_related` — Related categories.
- `get_category_series` — Series in a category.
- `get_category_tags` — Tags for a category.
- `get_category_related_tags` — Related tags.

### Releases

Releases are publications (e.g., GDP release `53`).

- `get_releases` — List all releases.
- `get_releases_dates` — Release dates.
- `get_release` — Release by ID.
- `get_release_dates` — Dates for a specific release.
- `get_release_sources` — Sources for a release.
- `get_release_series` — Series in a release.
- `get_release_tags` — Tags for a release.
- `get_release_related_tags` — Related tags for a release.
- `get_release_tables` — Release table structure.

### Tags

Tags are keywords (e.g., `gdp`, `inflation`, `business`).

- `get_tags` — Browse/search all tags. Pass `tag_names` (comma-separated) to filter, `search_text` for keyword search, `limit` to cap results.
- `get_related_tags` — Find tags related to given tag names.
- `get_tags_series` — Get series matching given tags.

### Search

Search across all FRED series. Results auto-paginate — set `limit` to control count.

- `search_series` — Search by keyword (`search_text`). Supports `order_by` (popularity, title, etc.), `sort_order`, `tag_names` filter, `exclude_tag_names`, `filter_variable`/`filter_value`.
- `get_release_series` — Series in a release.
- `get_category_series` — Series in a category.

### Sources

Data providers (e.g., Federal Reserve = `1`, BLS = `2`).

- `get_sources` — List all sources. Supports `limit` and `sort_order`.
- `get_source` — Source by ID.
- `get_source_releases` — Releases for a source.

### GeoFRED

Regional economic data with map support.

- `get_series_group` — GeoFRED metadata for a series. Returns region type, units, date range.
- `get_series_data` — Map data for a series (optionally filtered by date).
- `get_regional_data` — Regional data by series group. Requires `series_group`. Optional: `region_type` (state, msa, county), `season` (SA, NSA), `units`, `frequency`.

## Behavioral Rules

### Parameter Formatting

- **IDs are strings** — Category, release, and source IDs are passed as strings (e.g., `"1"`, `"53"`) even though they represent integers.
- **Tag names** — Pass multiple tags as comma-separated strings (e.g., `"gdp,inflation"`).
- **Dates** — Always use YYYY-MM-DD format.
- **Limit** — Observation tools use `limit` as a number. Search and tag tools use `limit` as a string.

### Common Workflows

**Find a series**: `search_series` with keywords → pick from results → `get_series_info` for metadata → `get_series_observations` for data.

**Compare revisions**: `get_series_all_releases` with `realtime_start` / `realtime_end` filters.

**Explore a category**: Start at root (`get_category_children` with `"0"`) → drill down → `get_category_series` → pull observations.

**Get the latest data point**: `get_series_observations` with `sort_order="desc"` and `limit="1"`.

### Known Reference IDs

These IDs are stable and can be used directly:

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
| Tag      | `gdp`       | GDP-related series                        |
| Tag      | `inflation` | Inflation-related series                  |
