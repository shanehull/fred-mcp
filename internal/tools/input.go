package tools

// Input structs for tool handlers. Fields without omitempty are required in
// the inferred JSON schema; shared option structs are embedded and flattened.

// ObservationOptions holds optional parameters shared by series observation tools.
type ObservationOptions struct {
	ObservationStart  string `json:"observation_start,omitempty" jsonschema:"Start of observation period in YYYY-MM-DD format."`
	ObservationEnd    string `json:"observation_end,omitempty" jsonschema:"End of observation period in YYYY-MM-DD format."`
	RealtimeStart     string `json:"realtime_start,omitempty" jsonschema:"Start of real-time period in YYYY-MM-DD format."`
	RealtimeEnd       string `json:"realtime_end,omitempty" jsonschema:"End of real-time period in YYYY-MM-DD format."`
	Units             string `json:"units,omitempty" jsonschema:"Units: lin, chg, ch1, pch, pc1, pca, cch, cca, log."`
	Frequency         string `json:"frequency,omitempty" jsonschema:"Frequency: d, w, bw, m, q, sa, a."`
	AggregationMethod string `json:"aggregation_method,omitempty" jsonschema:"Aggregation: avg, sum, eop."`
	OutputType        int    `json:"output_type,omitempty" jsonschema:"Output type (1-4)."`
	VintageDates      string `json:"vintage_dates,omitempty" jsonschema:"Comma-separated vintage dates in YYYY-MM-DD format."`
	SortOrder         string `json:"sort_order,omitempty" jsonschema:"Sort order: 'asc' or 'desc'."`
	Limit             int    `json:"limit,omitempty" jsonschema:"Maximum number of observations."`
	Offset            int    `json:"offset,omitempty" jsonschema:"Offset for pagination."`
}

// TagOptions holds optional parameters shared by tag-related tools.
type TagOptions struct {
	TagGroupID      string `json:"tag_group_id,omitempty" jsonschema:"Filter by tag group ID."`
	SearchText      string `json:"search_text,omitempty" jsonschema:"Search text for tags."`
	Limit           int    `json:"limit,omitempty" jsonschema:"Maximum number of tags."`
	OrderBy         string `json:"order_by,omitempty" jsonschema:"Field to order by."`
	SortOrder       string `json:"sort_order,omitempty" jsonschema:"Sort order: 'asc' or 'desc'."`
	ExcludeTagNames string `json:"exclude_tag_names,omitempty" jsonschema:"Comma-separated tags to exclude."`
}

// SearchOptions holds optional parameters shared by series search tools.
type SearchOptions struct {
	SearchType      string `json:"search_type,omitempty" jsonschema:"'full_text' or 'series_id'."`
	OrderBy         string `json:"order_by,omitempty" jsonschema:"Field to order by."`
	SortOrder       string `json:"sort_order,omitempty" jsonschema:"'asc' or 'desc'."`
	FilterVariable  string `json:"filter_variable,omitempty" jsonschema:"Filter variable name."`
	FilterValue     string `json:"filter_value,omitempty" jsonschema:"Filter variable value."`
	Limit           int    `json:"limit,omitempty" jsonschema:"Maximum number of results."`
	ExcludeTagNames string `json:"exclude_tag_names,omitempty" jsonschema:"Comma-separated tag names to exclude."`
}

// ReleaseListOptions holds optional parameters for release/source list tools.
type ReleaseListOptions struct {
	Limit     int    `json:"limit,omitempty" jsonschema:"Maximum number of results."`
	SortOrder string `json:"sort_order,omitempty" jsonschema:"'asc' or 'desc'."`
}

// ReleaseDateOptions holds optional parameters for release date tools.
type ReleaseDateOptions struct {
	Limit                         int    `json:"limit,omitempty" jsonschema:"Maximum number of results."`
	SortOrder                     string `json:"sort_order,omitempty" jsonschema:"'asc' or 'desc'."`
	IncludeReleaseDatesWithNoData string `json:"include_release_dates_with_no_data,omitempty" jsonschema:"Set to 'true' to include releases with no data."`
}

// TableOptions holds optional parameters for the release tables tool.
type TableOptions struct {
	ElementID               int    `json:"element_id,omitempty" jsonschema:"Element ID to filter by."`
	IncludeObservationValues string `json:"include_observation_values,omitempty" jsonschema:"Set to 'true' to include observations."`
	ObservationDate         string `json:"observation_date,omitempty" jsonschema:"Observation date in YYYY-MM-DD format."`
}

// UpdateOptions holds optional parameters for the series updates tool.
type UpdateOptions struct {
	StartTime   string `json:"start_time,omitempty" jsonschema:"Start time for updates filter."`
	EndTime     string `json:"end_time,omitempty" jsonschema:"End time for updates filter."`
	FilterValue string `json:"filter_value,omitempty" jsonschema:"Filter by 'macro', 'regional', or 'all'."`
	Limit       int    `json:"limit,omitempty" jsonschema:"Maximum number of updated series to return."`
}

// SourceOptions holds optional parameters for the get_sources tool.
type SourceOptions struct {
	Limit     int    `json:"limit,omitempty" jsonschema:"Maximum number of sources."`
	SortOrder string `json:"sort_order,omitempty" jsonschema:"'asc' or 'desc'."`
}

// MapDataOptions holds optional parameters for the GeoFRED series data tool.
type MapDataOptions struct {
	Date      string `json:"date,omitempty" jsonschema:"Map date in YYYY-MM-DD format."`
	StartDate string `json:"start_date,omitempty" jsonschema:"Map start date in YYYY-MM-DD format."`
}

// SeriesIDInput is the input for tools keyed only by a series ID.
type SeriesIDInput struct {
	SeriesID string `json:"series_id" jsonschema:"The series ID (e.g., 'DGS20', 'GDP')."`
}

// SeriesObservationInput is the input for observation tools.
type SeriesObservationInput struct {
	SeriesID string `json:"series_id" jsonschema:"The series ID (e.g., 'DGS20', 'GDP')."`
	ObservationOptions
}

// SeriesAsOfInput is the input for the get_series_as_of tool.
type SeriesAsOfInput struct {
	SeriesID string `json:"series_id" jsonschema:"The series ID (e.g., 'DGS20', 'GDP')."`
	AsOfDate string `json:"as_of_date" jsonschema:"Date in YYYY-MM-DD format."`
	ObservationOptions
}

// SeriesIDTagInput is the input for series tag tools.
type SeriesIDTagInput struct {
	SeriesID string `json:"series_id" jsonschema:"The series ID (e.g., 'DGS20', 'GDP')."`
	TagOptions
}

// SearchSeriesTagsInput is the input for the search_series_tags tool.
type SearchSeriesTagsInput struct {
	SeriesSearchText string `json:"series_search_text" jsonschema:"Text to search for in series tags."`
	TagOptions
}

// SearchSeriesRelatedTagsInput is the input for the search_series_related_tags tool.
type SearchSeriesRelatedTagsInput struct {
	SeriesSearchText string `json:"series_search_text" jsonschema:"Text to search for."`
	TagNames         string `json:"tag_names,omitempty" jsonschema:"Comma-separated list of tag names."`
	TagOptions
}

// CategoryIDInput is the input for tools keyed only by a category ID.
type CategoryIDInput struct {
	CategoryID *int `json:"category_id" jsonschema:"The category ID (integer)."`
}

// CategoryIDTagInput is the input for category tag tools.
type CategoryIDTagInput struct {
	CategoryID *int `json:"category_id" jsonschema:"The category ID (integer)."`
	TagOptions
}

// CategoryRelatedTagsInput is the input for the get_category_related_tags tool.
type CategoryRelatedTagsInput struct {
	CategoryID *int `json:"category_id" jsonschema:"The category ID (integer)."`
	TagNames   string `json:"tag_names,omitempty" jsonschema:"Comma-separated list of tag names."`
	TagOptions
}

// ReleaseIDInput is the input for tools keyed only by a release ID.
type ReleaseIDInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
}

// ReleaseIDDatesInput is the input for the get_release_dates tool.
type ReleaseIDDatesInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
	ReleaseDateOptions
}

// ReleaseIDTagInput is the input for release tag tools.
type ReleaseIDTagInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
	TagOptions
}

// ReleaseRelatedTagsInput is the input for the get_release_related_tags tool.
type ReleaseRelatedTagsInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
	TagNames  string `json:"tag_names,omitempty" jsonschema:"Comma-separated list of tag names."`
	TagOptions
}

// ReleaseTablesInput is the input for the get_release_tables tool.
type ReleaseTablesInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
	TableOptions
}

// SourceIDInput is the input for tools keyed only by a source ID.
type SourceIDInput struct {
	SourceID *int `json:"source_id" jsonschema:"The source ID (integer)."`
}

// SourceIDReleasesInput is the input for the get_source_releases tool.
type SourceIDReleasesInput struct {
	SourceID *int `json:"source_id" jsonschema:"The source ID (integer)."`
	ReleaseListOptions
}

// TagsInput is the input for the get_tags tool.
type TagsInput struct {
	TagNames string `json:"tag_names,omitempty" jsonschema:"Comma-separated list of tag names to filter by."`
	TagOptions
}

// RelatedTagsInput is the input for the get_related_tags tool.
type RelatedTagsInput struct {
	TagNames string `json:"tag_names" jsonschema:"Comma-separated list of tag names."`
	TagOptions
}

// TagsSeriesInput is the input for the get_tags_series tool.
type TagsSeriesInput struct {
	TagNames string `json:"tag_names" jsonschema:"Comma-separated list of tag names."`
	SearchOptions
}

// SearchSeriesInput is the input for the search_series tool.
type SearchSeriesInput struct {
	SearchText string `json:"search_text" jsonschema:"Text to search for."`
	TagNames   string `json:"tag_names,omitempty" jsonschema:"Comma-separated tag names to filter by."`
	SearchOptions
}

// ReleaseSeriesInput is the input for the get_release_series tool.
type ReleaseSeriesInput struct {
	ReleaseID *int `json:"release_id" jsonschema:"The release ID (integer)."`
	TagNames  string `json:"tag_names,omitempty" jsonschema:"Comma-separated tag names to filter by."`
	SearchOptions
}

// CategorySeriesInput is the input for the get_category_series tool.
type CategorySeriesInput struct {
	CategoryID *int `json:"category_id" jsonschema:"The category ID (integer)."`
	TagNames   string `json:"tag_names,omitempty" jsonschema:"Comma-separated tag names to filter by."`
	SearchOptions
}

// SeriesDataInput is the input for the get_series_data tool.
type SeriesDataInput struct {
	SeriesID string `json:"series_id" jsonschema:"The GeoFRED series ID (e.g., 'WIPCPI')."`
	MapDataOptions
}

// SeriesGroupInput is the input for the get_series_group tool.
type SeriesGroupInput struct {
	SeriesID string `json:"series_id" jsonschema:"The GeoFRED series ID (e.g., 'WIPCPI')."`
}

// RegionalDataInput is the input for the get_regional_data tool.
type RegionalDataInput struct {
	SeriesGroup    string `json:"series_group" jsonschema:"The series group ID (required)."`
	Date           string `json:"date" jsonschema:"Map date in YYYY-MM-DD format."`
	Frequency      string `json:"frequency" jsonschema:"Frequency aggregation: d, w, bw, m, q, sa, a."`
	Units          string `json:"units" jsonschema:"Units: lin, chg, ch1, pch, pc1, pca, cch, cca, log."`
	Season         string `json:"season" jsonschema:"Seasonality: SA, NSA, SSA, SAAR, NSAAR."`
	RegionType     string `json:"region_type" jsonschema:"Region type: bea, msa, frb, necta, state, country, county, censusregion."`
	Transformation string `json:"transformation,omitempty" jsonschema:"Transformation: lin, chg, ch1, pch, pc1, pca, cch, cca, log."`
}
