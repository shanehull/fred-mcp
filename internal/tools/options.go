package tools

import (
	"strings"
	"time"

	"github.com/shanehull/go-fred"
)

func buildObservationOptions(in ObservationOptions) []fred.ObservationOption {
	var opts []fred.ObservationOption

	if in.ObservationStart != "" {
		if t, err := time.Parse("2006-01-02", in.ObservationStart); err == nil {
			opts = append(opts, fred.WithObservationStart(t))
		}
	}
	if in.ObservationEnd != "" {
		if t, err := time.Parse("2006-01-02", in.ObservationEnd); err == nil {
			opts = append(opts, fred.WithObservationEnd(t))
		}
	}
	if in.RealtimeStart != "" {
		if t, err := time.Parse("2006-01-02", in.RealtimeStart); err == nil {
			opts = append(opts, fred.WithRealtimeStart(t))
		}
	}
	if in.RealtimeEnd != "" {
		if t, err := time.Parse("2006-01-02", in.RealtimeEnd); err == nil {
			opts = append(opts, fred.WithRealtimeEnd(t))
		}
	}
	if in.Units != "" {
		opts = append(opts, fred.WithUnits(in.Units))
	}
	if in.Frequency != "" {
		opts = append(opts, fred.WithFrequency(in.Frequency))
	}
	if in.AggregationMethod != "" {
		opts = append(opts, fred.WithAggregationMethod(in.AggregationMethod))
	}
	if in.OutputType != 0 {
		opts = append(opts, fred.WithOutputType(in.OutputType))
	}
	if tags := parseStringList(in.VintageDates); len(tags) > 0 {
		opts = append(opts, fred.WithVintageDates(tags...))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithObservationSortOrder(fred.SortOrder(in.SortOrder)))
	}
	if in.Limit != 0 {
		opts = append(opts, fred.WithObservationLimit(in.Limit))
	}
	if in.Offset != 0 {
		opts = append(opts, fred.WithObservationOffset(in.Offset))
	}

	return opts
}

func buildSearchOptions(in SearchOptions) []fred.SearchOption {
	var opts []fred.SearchOption

	if in.SearchType != "" {
		opts = append(opts, fred.WithSearchType(in.SearchType))
	}
	if in.OrderBy != "" {
		opts = append(opts, fred.WithOrderBy(fred.OrderBy(in.OrderBy)))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithSortOrder(fred.SortOrder(in.SortOrder)))
	}
	if in.FilterVariable != "" {
		opts = append(opts, fred.WithFilter(in.FilterVariable, in.FilterValue))
	}
	if in.Limit > 0 {
		opts = append(opts, fred.WithLimit(in.Limit))
	}
	if tags := parseStringList(in.ExcludeTagNames); len(tags) > 0 {
		opts = append(opts, fred.WithExcludeTags(tags...))
	}

	return opts
}

func buildTagOptions(in TagOptions) []fred.TagOption {
	var opts []fred.TagOption

	if in.TagGroupID != "" {
		opts = append(opts, fred.WithTagGroupID(in.TagGroupID))
	}
	if in.SearchText != "" {
		opts = append(opts, fred.WithTagSearchText(in.SearchText))
	}
	if in.Limit != 0 {
		opts = append(opts, fred.WithTagLimit(in.Limit))
	}
	if in.OrderBy != "" {
		opts = append(opts, fred.WithTagOrderBy(fred.OrderBy(in.OrderBy)))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithTagSortOrder(fred.SortOrder(in.SortOrder)))
	}
	if tags := parseStringList(in.ExcludeTagNames); len(tags) > 0 {
		opts = append(opts, fred.WithTagSetExclude(tags...))
	}

	return opts
}

func buildReleaseListOptions(in ReleaseListOptions) []fred.ReleaseListOption {
	var opts []fred.ReleaseListOption

	if in.Limit != 0 {
		opts = append(opts, fred.WithReleaseLimit(in.Limit))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithReleaseSortOrder(fred.SortOrder(in.SortOrder)))
	}

	return opts
}

func buildReleaseDateOptions(in ReleaseDateOptions) []fred.ReleaseDateOption {
	var opts []fred.ReleaseDateOption

	if in.Limit != 0 {
		opts = append(opts, fred.WithReleaseDateLimit(in.Limit))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithReleaseDateSortOrder(fred.SortOrder(in.SortOrder)))
	}
	if in.IncludeReleaseDatesWithNoData == "true" {
		opts = append(opts, fred.WithIncludeNoData(true))
	}

	return opts
}

func buildTableOptions(in TableOptions) []fred.TableOption {
	var opts []fred.TableOption

	if in.ElementID != 0 {
		opts = append(opts, fred.WithTableElementID(in.ElementID))
	}
	if in.IncludeObservationValues == "true" {
		opts = append(opts, fred.WithIncludeObservationValues(true))
	}
	if in.ObservationDate != "" {
		opts = append(opts, fred.WithObservationDate(in.ObservationDate))
	}

	return opts
}

func buildUpdateOptions(in UpdateOptions) []fred.UpdateOption {
	var opts []fred.UpdateOption

	if in.StartTime != "" {
		opts = append(opts, fred.WithStartTime(in.StartTime))
	}
	if in.EndTime != "" {
		opts = append(opts, fred.WithEndTime(in.EndTime))
	}
	if in.FilterValue != "" {
		opts = append(opts, fred.WithFilterValue(in.FilterValue))
	}
	if in.Limit != 0 {
		opts = append(opts, fred.WithUpdateLimit(in.Limit))
	}

	return opts
}

func buildSourceOptions(in SourceOptions) []fred.SourceOption {
	var opts []fred.SourceOption

	if in.Limit != 0 {
		opts = append(opts, fred.WithSourceLimit(in.Limit))
	}
	if in.SortOrder != "" {
		opts = append(opts, fred.WithSourceSortOrder(fred.SortOrder(in.SortOrder)))
	}

	return opts
}

func buildMapDataOptions(in MapDataOptions) []fred.MapDataOption {
	var opts []fred.MapDataOption

	if in.Date != "" {
		opts = append(opts, fred.WithMapDate(in.Date))
	}
	if in.StartDate != "" {
		opts = append(opts, fred.WithMapStartDate(in.StartDate))
	}

	return opts
}

func buildRegionalDataOptions(in RegionalDataInput) []fred.RegionalDataOption {
	var opts []fred.RegionalDataOption

	if in.SeriesGroup != "" {
		opts = append(opts, fred.WithSeriesGroup(in.SeriesGroup))
	}
	if in.RegionType != "" {
		opts = append(opts, fred.WithRegionType(in.RegionType))
	}
	if in.Date != "" {
		opts = append(opts, fred.WithRegionalDate(in.Date))
	}
	if in.Season != "" {
		opts = append(opts, fred.WithSeason(in.Season))
	}
	if in.Units != "" {
		opts = append(opts, fred.WithMapUnits(in.Units))
	}
	if in.Transformation != "" {
		opts = append(opts, fred.WithTransformation(in.Transformation))
	}
	if in.Frequency != "" {
		opts = append(opts, fred.WithRegionalFrequency(in.Frequency))
	}

	return opts
}

func parseStringList(raw string) []string {
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
