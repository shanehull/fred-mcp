package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func buildObservationOptions(req mcp.CallToolRequest) []fred.ObservationOption {
	var opts []fred.ObservationOption

	if v := req.GetString("observation_start", ""); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts = append(opts, fred.WithObservationStart(t))
		}
	}
	if v := req.GetString("observation_end", ""); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts = append(opts, fred.WithObservationEnd(t))
		}
	}
	if v := req.GetString("realtime_start", ""); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts = append(opts, fred.WithRealtimeStart(t))
		}
	}
	if v := req.GetString("realtime_end", ""); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts = append(opts, fred.WithRealtimeEnd(t))
		}
	}
	if v := req.GetString("units", ""); v != "" {
		opts = append(opts, fred.WithUnits(v))
	}
	if v := req.GetString("frequency", ""); v != "" {
		opts = append(opts, fred.WithFrequency(v))
	}
	if v := req.GetString("aggregation_method", ""); v != "" {
		opts = append(opts, fred.WithAggregationMethod(v))
	}
	if v := req.GetFloat("output_type", 0); v != 0 {
		opts = append(opts, fred.WithOutputType(int(v)))
	}
	if v := parseStringList(req, "vintage_dates"); len(v) > 0 {
		opts = append(opts, fred.WithVintageDates(v...))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithObservationSortOrder(fred.SortOrder(v)))
	}
	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithObservationLimit(int(v)))
	}
	if v := req.GetFloat("offset", 0); v != 0 {
		opts = append(opts, fred.WithObservationOffset(int(v)))
	}

	return opts
}

func buildSearchOptions(req mcp.CallToolRequest) []fred.SearchOption {
	var opts []fred.SearchOption

	if v := req.GetString("search_type", ""); v != "" {
		opts = append(opts, fred.WithSearchType(v))
	}
	if v := req.GetString("order_by", ""); v != "" {
		opts = append(opts, fred.WithOrderBy(fred.OrderBy(v)))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithSortOrder(fred.SortOrder(v)))
	}
	if filterVar := req.GetString("filter_variable", ""); filterVar != "" {
		filterVal := req.GetString("filter_value", "")
		opts = append(opts, fred.WithFilter(filterVar, filterVal))
	}
	if v := req.GetFloat("limit", 0); v > 0 {
		opts = append(opts, fred.WithLimit(int(v)))
	}
	if tags := parseStringList(req, "tag_names"); len(tags) > 0 {
		opts = append(opts, fred.WithTagNames(tags...))
	}
	if tags := parseStringList(req, "exclude_tag_names"); len(tags) > 0 {
		opts = append(opts, fred.WithExcludeTags(tags...))
	}

	return opts
}

func buildTagOptions(req mcp.CallToolRequest) []fred.TagOption {
	var opts []fred.TagOption

	if v := req.GetString("tag_group_id", ""); v != "" {
		opts = append(opts, fred.WithTagGroupID(v))
	}
	if v := req.GetString("search_text", ""); v != "" {
		opts = append(opts, fred.WithTagSearchText(v))
	}
	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithTagLimit(int(v)))
	}
	if v := req.GetString("order_by", ""); v != "" {
		opts = append(opts, fred.WithTagOrderBy(fred.OrderBy(v)))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithTagSortOrder(fred.SortOrder(v)))
	}
	if tags := parseStringList(req, "tag_names"); len(tags) > 0 {
		opts = append(opts, fred.WithTagSetNames(tags...))
	}
	if tags := parseStringList(req, "exclude_tag_names"); len(tags) > 0 {
		opts = append(opts, fred.WithTagSetExclude(tags...))
	}

	return opts
}

func buildReleaseListOptions(req mcp.CallToolRequest) []fred.ReleaseListOption {
	var opts []fred.ReleaseListOption

	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithReleaseLimit(int(v)))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithReleaseSortOrder(fred.SortOrder(v)))
	}

	return opts
}

func buildReleaseDateOptions(req mcp.CallToolRequest) []fred.ReleaseDateOption {
	var opts []fred.ReleaseDateOption

	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithReleaseDateLimit(int(v)))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithReleaseDateSortOrder(fred.SortOrder(v)))
	}
	if v := req.GetString("include_release_dates_with_no_data", ""); v == "true" {
		opts = append(opts, fred.WithIncludeNoData(true))
	}

	return opts
}

func buildTableOptions(req mcp.CallToolRequest) []fred.TableOption {
	var opts []fred.TableOption

	if v := req.GetFloat("element_id", 0); v != 0 {
		opts = append(opts, fred.WithTableElementID(int(v)))
	}
	if v := req.GetString("include_observation_values", ""); v == "true" {
		opts = append(opts, fred.WithIncludeObservationValues(true))
	}
	if v := req.GetString("observation_date", ""); v != "" {
		opts = append(opts, fred.WithObservationDate(v))
	}

	return opts
}

func buildUpdateOptions(req mcp.CallToolRequest) []fred.UpdateOption {
	var opts []fred.UpdateOption

	if v := req.GetString("start_time", ""); v != "" {
		opts = append(opts, fred.WithStartTime(v))
	}
	if v := req.GetString("end_time", ""); v != "" {
		opts = append(opts, fred.WithEndTime(v))
	}
	if v := req.GetString("filter_value", ""); v != "" {
		opts = append(opts, fred.WithFilterValue(v))
	}
	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithUpdateLimit(int(v)))
	}

	return opts
}

func buildSourceOptions(req mcp.CallToolRequest) []fred.SourceOption {
	var opts []fred.SourceOption

	if v := req.GetFloat("limit", 0); v != 0 {
		opts = append(opts, fred.WithSourceLimit(int(v)))
	}
	if v := req.GetString("sort_order", ""); v != "" {
		opts = append(opts, fred.WithSourceSortOrder(fred.SortOrder(v)))
	}

	return opts
}

func buildMapDataOptions(req mcp.CallToolRequest) []fred.MapDataOption {
	var opts []fred.MapDataOption

	if v := req.GetString("date", ""); v != "" {
		opts = append(opts, fred.WithMapDate(v))
	}
	if v := req.GetString("start_date", ""); v != "" {
		opts = append(opts, fred.WithMapStartDate(v))
	}

	return opts
}

func buildRegionalDataOptions(req mcp.CallToolRequest) []fred.RegionalDataOption {
	var opts []fred.RegionalDataOption

	if v := req.GetString("series_group", ""); v != "" {
		opts = append(opts, fred.WithSeriesGroup(v))
	}
	if v := req.GetString("region_type", ""); v != "" {
		opts = append(opts, fred.WithRegionType(v))
	}
	if v := req.GetString("date", ""); v != "" {
		opts = append(opts, fred.WithRegionalDate(v))
	}
	if v := req.GetString("season", ""); v != "" {
		opts = append(opts, fred.WithSeason(v))
	}
	if v := req.GetString("units", ""); v != "" {
		opts = append(opts, fred.WithMapUnits(v))
	}
	if v := req.GetString("transformation", ""); v != "" {
		opts = append(opts, fred.WithTransformation(v))
	}
	if v := req.GetString("frequency", ""); v != "" {
		opts = append(opts, fred.WithRegionalFrequency(v))
	}

	return opts
}

func parseStringList(req mcp.CallToolRequest, key string) []string {
	raw := req.GetString(key, "")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func GetInt(req mcp.CallToolRequest, key string) (int, error) {
	raw := req.GetString(key, "")
	if raw == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	return strconv.Atoi(raw)
}

func GetIntPtr(req mcp.CallToolRequest, key string) (int, bool) {
	raw := req.GetString(key, "")
	if raw == "" {
		return 0, false
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}
