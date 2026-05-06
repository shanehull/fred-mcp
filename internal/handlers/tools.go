package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shanehull/fred-mcp/internal/tools"
	"github.com/shanehull/go-fred"
)

func RegisterTools(s *server.MCPServer, client *fred.Client) {
	registerSeriesTools(s, client)
	registerCategoryTools(s, client)
	registerReleaseTools(s, client)
	registerSourceTools(s, client)
	registerTagTools(s, client)
	registerSearchTools(s, client)
	registerGeoFREDTools(s, client)
}

func registerSeriesTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_series_info",
		mcp.WithDescription("Get an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID (e.g., 'DGS20', 'GDP').")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesInfo))

	obsOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the observations or data values for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
	}
	obsOpts = append(obsOpts, withObservationOptions()...)
	s.AddTool(mcp.NewTool("get_series_observations", obsOpts...),
		toolHandler(client, tools.HandleGetSeriesObservations))

	allOpts := []mcp.ToolOption{
		mcp.WithDescription("Get all releases for an economic data series (all vintages)."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
	}
	allOpts = append(allOpts, withObservationOptions()...)
	s.AddTool(mcp.NewTool("get_series_all_releases", allOpts...),
		toolHandler(client, tools.HandleGetSeriesAllReleases))

	firstOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the first release of observations for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
	}
	firstOpts = append(firstOpts, withObservationOptions()...)
	s.AddTool(mcp.NewTool("get_series_first_release", firstOpts...),
		toolHandler(client, tools.HandleGetSeriesFirstRelease))

	asofOpts := []mcp.ToolOption{
		mcp.WithDescription("Get observations for a series as of a specific date."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
		mcp.WithString("as_of_date", mcp.Required(), mcp.Description("Date in YYYY-MM-DD format.")),
	}
	asofOpts = append(asofOpts, withObservationOptions()...)
	s.AddTool(mcp.NewTool("get_series_as_of", asofOpts...),
		toolHandler(client, tools.HandleGetSeriesAsOf))

	vintageOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the vintage dates for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
	}
	vintageOpts = append(vintageOpts, withObservationOptions()...)
	s.AddTool(mcp.NewTool("get_series_vintage_dates", vintageOpts...),
		toolHandler(client, tools.HandleGetSeriesVintageDates))

	s.AddTool(mcp.NewTool("get_series_categories",
		mcp.WithDescription("Get the categories for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesCategories))

	s.AddTool(mcp.NewTool("get_series_release",
		mcp.WithDescription("Get the release for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesRelease))

	seriesTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the tags for an economic data series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
	}
	seriesTagsOpts = append(seriesTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("get_series_tags", seriesTagsOpts...),
		toolHandler(client, tools.HandleGetSeriesTags))

	searchTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Search series tags by text."),
		mcp.WithString("series_search_text", mcp.Required(), mcp.Description("Text to search for in series tags.")),
	}
	searchTagsOpts = append(searchTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("search_series_tags", searchTagsOpts...),
		toolHandler(client, tools.HandleSearchSeriesTags))

	relatedTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Search for related series tags by text."),
		mcp.WithString("series_search_text", mcp.Required(), mcp.Description("Text to search for.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated list of tag names.")),
	}
	relatedTagsOpts = append(relatedTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("search_series_related_tags", relatedTagsOpts...),
		toolHandler(client, tools.HandleSearchSeriesRelatedTags))

	s.AddTool(mcp.NewTool("get_series_updates",
		mcp.WithDescription("Get economic data series that were updated recently."),
		mcp.WithString("start_time", mcp.Description("Start time for updates filter.")),
		mcp.WithString("end_time", mcp.Description("End time for updates filter.")),
		mcp.WithString("filter_value", mcp.Description("Filter by 'macro', 'regional', or 'all'.")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of updated series to return.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesUpdates))
}

func registerCategoryTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_category",
		mcp.WithDescription("Get a FRED category by ID."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID (integer).")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetCategory))

	s.AddTool(mcp.NewTool("get_category_children",
		mcp.WithDescription("Get the child categories for a specified parent category."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The parent category ID.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetCategoryChildren))

	s.AddTool(mcp.NewTool("get_category_related",
		mcp.WithDescription("Get related categories for a category."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetCategoryRelated))

	catTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the tags for a category."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID.")),
	}
	catTagsOpts = append(catTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("get_category_tags", catTagsOpts...),
		toolHandler(client, tools.HandleGetCategoryTags))

	catRelatedTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the related tags for a category."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated list of tag names.")),
	}
	catRelatedTagsOpts = append(catRelatedTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("get_category_related_tags", catRelatedTagsOpts...),
		toolHandler(client, tools.HandleGetCategoryRelatedTags))
}

func registerReleaseTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_releases",
		mcp.WithDescription("Get all releases of economic data."),
		mcp.WithString("limit", mcp.Description("Maximum number of releases to return.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleases))

	s.AddTool(mcp.NewTool("get_releases_dates",
		mcp.WithDescription("Get release dates for all releases of economic data."),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("include_release_dates_with_no_data", mcp.Description("Set to 'true' to include releases with no data.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleasesDates))

	s.AddTool(mcp.NewTool("get_release",
		mcp.WithDescription("Get a specific release of economic data."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetRelease))

	s.AddTool(mcp.NewTool("get_release_dates",
		mcp.WithDescription("Get release dates for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("include_release_dates_with_no_data", mcp.Description("Set to 'true' to include releases with no data.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleaseDates))

	s.AddTool(mcp.NewTool("get_release_sources",
		mcp.WithDescription("Get the sources for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleaseSources))

	relTagsOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the tags for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
	}
	relTagsOpts = append(relTagsOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("get_release_tags", relTagsOpts...),
		toolHandler(client, tools.HandleGetReleaseTags))

	relRelatedOpts := []mcp.ToolOption{
		mcp.WithDescription("Get the related tags for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated list of tag names.")),
	}
	relRelatedOpts = append(relRelatedOpts, withTagOptions()...)
	s.AddTool(mcp.NewTool("get_release_related_tags", relRelatedOpts...),
		toolHandler(client, tools.HandleGetReleaseRelatedTags))

	s.AddTool(mcp.NewTool("get_release_tables",
		mcp.WithDescription("Get the release tables for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithString("element_id", mcp.Description("Element ID to filter by.")),
		mcp.WithString("include_observation_values", mcp.Description("Set to 'true' to include observations.")),
		mcp.WithString("observation_date", mcp.Description("Observation date in YYYY-MM-DD format.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleaseTables))
}

func registerSourceTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_sources",
		mcp.WithDescription("Get all sources of economic data."),
		mcp.WithString("limit", mcp.Description("Maximum number of sources.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSources))

	s.AddTool(mcp.NewTool("get_source",
		mcp.WithDescription("Get a specific source of economic data."),
		mcp.WithString("source_id", mcp.Required(), mcp.Description("The source ID (integer).")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSource))

	s.AddTool(mcp.NewTool("get_source_releases",
		mcp.WithDescription("Get the releases for a specific source."),
		mcp.WithString("source_id", mcp.Required(), mcp.Description("The source ID (integer).")),
		mcp.WithString("limit", mcp.Description("Maximum number of releases.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSourceReleases))
}

func registerTagTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_tags",
		mcp.WithDescription("Get FRED tags."),
		mcp.WithString("tag_group_id", mcp.Description("Filter by tag group ID.")),
		mcp.WithString("search_text", mcp.Description("Search text for tags.")),
		mcp.WithString("limit", mcp.Description("Maximum number of tags.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated list of tag names to filter by.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated list of tag names to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetTags))

	s.AddTool(mcp.NewTool("get_related_tags",
		mcp.WithDescription("Get related FRED tags for one or more tag names."),
		mcp.WithString("tag_names", mcp.Required(), mcp.Description("Comma-separated list of tag names.")),
		mcp.WithString("tag_group_id", mcp.Description("Filter by tag group ID.")),
		mcp.WithString("search_text", mcp.Description("Search text for tags.")),
		mcp.WithString("limit", mcp.Description("Maximum number of tags.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated list of tag names to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetRelatedTags))

	s.AddTool(mcp.NewTool("get_tags_series",
		mcp.WithDescription("Get the series matching one or more tags."),
		mcp.WithString("tag_names", mcp.Required(), mcp.Description("Comma-separated list of tag names.")),
		mcp.WithString("search_type", mcp.Description("'full_text' or 'series_id'.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("filter_variable", mcp.Description("Filter variable name.")),
		mcp.WithString("filter_value", mcp.Description("Filter variable value.")),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated tags to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetTagsSeries))
}

func registerSearchTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("search_series",
		mcp.WithDescription("Search for economic data series by text."),
		mcp.WithString("search_text", mcp.Required(), mcp.Description("Text to search for.")),
		mcp.WithString("search_type", mcp.Description("'full_text' or 'series_id'.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("filter_variable", mcp.Description("Filter variable name.")),
		mcp.WithString("filter_value", mcp.Description("Filter variable value.")),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated tag names to filter by.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated tag names to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleSearchSeries))

	s.AddTool(mcp.NewTool("get_release_series",
		mcp.WithDescription("Get the series for a specific release."),
		mcp.WithString("release_id", mcp.Required(), mcp.Description("The release ID (integer).")),
		mcp.WithString("search_type", mcp.Description("'full_text' or 'series_id'.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("filter_variable", mcp.Description("Filter variable name.")),
		mcp.WithString("filter_value", mcp.Description("Filter variable value.")),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated tag names to filter by.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated tag names to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetReleaseSeries))

	s.AddTool(mcp.NewTool("get_category_series",
		mcp.WithDescription("Get the series for a specific category."),
		mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID (integer).")),
		mcp.WithString("search_type", mcp.Description("'full_text' or 'series_id'.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("'asc' or 'desc'.")),
		mcp.WithString("filter_variable", mcp.Description("Filter variable name.")),
		mcp.WithString("filter_value", mcp.Description("Filter variable value.")),
		mcp.WithString("limit", mcp.Description("Maximum number of results.")),
		mcp.WithString("tag_names", mcp.Description("Comma-separated tag names to filter by.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated tag names to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetCategorySeries))
}

func registerGeoFREDTools(s *server.MCPServer, client *fred.Client) {
	s.AddTool(mcp.NewTool("get_series_group",
		mcp.WithDescription("Get the series group for a GeoFRED series."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesGroup))

	s.AddTool(mcp.NewTool("get_series_data",
		mcp.WithDescription("Get GeoFRED series map data."),
		mcp.WithString("series_id", mcp.Required(), mcp.Description("The series ID.")),
		mcp.WithString("date", mcp.Description("Map date in YYYY-MM-DD format.")),
		mcp.WithString("start_date", mcp.Description("Map start date in YYYY-MM-DD format.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetSeriesData))

	s.AddTool(mcp.NewTool("get_regional_data",
		mcp.WithDescription("Get GeoFRED regional map data."),
		mcp.WithString("series_group", mcp.Required(), mcp.Description("The series group ID (required).")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Map date in YYYY-MM-DD format.")),
		mcp.WithString("frequency", mcp.Required(), mcp.Description("Frequency aggregation: d, w, bw, m, q, sa, a.")),
		mcp.WithString("units", mcp.Required(), mcp.Description("Units: lin, chg, ch1, pch, pc1, pca, cch, cca, log.")),
		mcp.WithString("season", mcp.Required(), mcp.Description("Seasonality: SA, NSA, SSA, SAAR, NSAAR.")),
		mcp.WithString("region_type", mcp.Required(), mcp.Description("Region type: bea, msa, frb, necta, state, country, county, censusregion.")),
		mcp.WithString("transformation", mcp.Description("Transformation: lin, chg, ch1, pch, pc1, pca, cch, cca, log.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	), toolHandler(client, tools.HandleGetRegionalData))
}

func toolHandler(client *fred.Client, fn func(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error)) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return fn(ctx, client, req)
	}
}

func withObservationOptions() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("observation_start", mcp.Description("Start of observation period in YYYY-MM-DD format.")),
		mcp.WithString("observation_end", mcp.Description("End of observation period in YYYY-MM-DD format.")),
		mcp.WithString("realtime_start", mcp.Description("Start of real-time period in YYYY-MM-DD format.")),
		mcp.WithString("realtime_end", mcp.Description("End of real-time period in YYYY-MM-DD format.")),
		mcp.WithString("units", mcp.Description("Units: lin, chg, ch1, pch, pc1, pca, cch, cca, log.")),
		mcp.WithString("frequency", mcp.Description("Frequency: d, w, bw, m, q, sa, a.")),
		mcp.WithString("aggregation_method", mcp.Description("Aggregation: avg, sum, eop.")),
		mcp.WithNumber("output_type", mcp.Description("Output type (1-4).")),
		mcp.WithString("vintage_dates", mcp.Description("Comma-separated vintage dates in YYYY-MM-DD format.")),
		mcp.WithString("sort_order", mcp.Description("Sort order: 'asc' or 'desc'.")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of observations.")),
		mcp.WithNumber("offset", mcp.Description("Offset for pagination.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	}
}

func withTagOptions() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithString("tag_group_id", mcp.Description("Filter by tag group ID.")),
		mcp.WithString("search_text", mcp.Description("Search text for tags.")),
		mcp.WithString("limit", mcp.Description("Maximum number of tags.")),
		mcp.WithString("order_by", mcp.Description("Field to order by.")),
		mcp.WithString("sort_order", mcp.Description("Sort order: 'asc' or 'desc'.")),
		mcp.WithString("exclude_tag_names", mcp.Description("Comma-separated tags to exclude.")),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	}
}
