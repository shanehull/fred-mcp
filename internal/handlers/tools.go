package handlers

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/fred-mcp/internal/tools"
	"github.com/shanehull/go-fred"
)

type toolFn[In any] func(ctx context.Context, client *fred.Client, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error)

func RegisterTools(s *mcp.Server, client *fred.Client) {
	registerSeriesTools(s, client)
	registerCategoryTools(s, client)
	registerReleaseTools(s, client)
	registerSourceTools(s, client)
	registerTagTools(s, client)
	registerSearchTools(s, client)
	registerGeoFREDTools(s, client)
}

func registerSeriesTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_series_info", "Get an economic data series."), toolHandler(client, tools.HandleGetSeriesInfo))
	mcp.AddTool(s, readOnlyTool("get_series_observations", "Get the observations or data values for an economic data series."), toolHandler(client, tools.HandleGetSeriesObservations))
	mcp.AddTool(s, readOnlyTool("get_series_all_releases", "Get all releases for an economic data series (all vintages)."), toolHandler(client, tools.HandleGetSeriesAllReleases))
	mcp.AddTool(s, readOnlyTool("get_series_first_release", "Get the first release of observations for an economic data series."), toolHandler(client, tools.HandleGetSeriesFirstRelease))
	mcp.AddTool(s, readOnlyTool("get_series_as_of", "Get observations for a series as of a specific date."), toolHandler(client, tools.HandleGetSeriesAsOf))
	mcp.AddTool(s, readOnlyTool("get_series_vintage_dates", "Get the vintage dates for an economic data series."), toolHandler(client, tools.HandleGetSeriesVintageDates))
	mcp.AddTool(s, readOnlyTool("get_series_categories", "Get the categories for an economic data series."), toolHandler(client, tools.HandleGetSeriesCategories))
	mcp.AddTool(s, readOnlyTool("get_series_release", "Get the release for an economic data series."), toolHandler(client, tools.HandleGetSeriesRelease))
	mcp.AddTool(s, readOnlyTool("get_series_tags", "Get the tags for an economic data series."), toolHandler(client, tools.HandleGetSeriesTags))
	mcp.AddTool(s, readOnlyTool("search_series_tags", "Search series tags by text."), toolHandler(client, tools.HandleSearchSeriesTags))
	mcp.AddTool(s, readOnlyTool("search_series_related_tags", "Search for related series tags by text."), toolHandler(client, tools.HandleSearchSeriesRelatedTags))
	mcp.AddTool(s, readOnlyTool("get_series_updates", "Get economic data series that were updated recently."), toolHandler(client, tools.HandleGetSeriesUpdates))
}

func registerCategoryTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_category", "Get a FRED category by ID."), toolHandler(client, tools.HandleGetCategory))
	mcp.AddTool(s, readOnlyTool("get_category_children", "Get the child categories for a specified parent category."), toolHandler(client, tools.HandleGetCategoryChildren))
	mcp.AddTool(s, readOnlyTool("get_category_related", "Get related categories for a category."), toolHandler(client, tools.HandleGetCategoryRelated))
	mcp.AddTool(s, readOnlyTool("get_category_tags", "Get the tags for a category."), toolHandler(client, tools.HandleGetCategoryTags))
	mcp.AddTool(s, readOnlyTool("get_category_related_tags", "Get the related tags for a category."), toolHandler(client, tools.HandleGetCategoryRelatedTags))
}

func registerReleaseTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_releases", "Get all releases of economic data."), toolHandler(client, tools.HandleGetReleases))
	mcp.AddTool(s, readOnlyTool("get_releases_dates", "Get release dates for all releases of economic data."), toolHandler(client, tools.HandleGetReleasesDates))
	mcp.AddTool(s, readOnlyTool("get_release", "Get a specific release of economic data."), toolHandler(client, tools.HandleGetRelease))
	mcp.AddTool(s, readOnlyTool("get_release_dates", "Get release dates for a specific release."), toolHandler(client, tools.HandleGetReleaseDates))
	mcp.AddTool(s, readOnlyTool("get_release_sources", "Get the sources for a specific release."), toolHandler(client, tools.HandleGetReleaseSources))
	mcp.AddTool(s, readOnlyTool("get_release_tags", "Get the tags for a specific release."), toolHandler(client, tools.HandleGetReleaseTags))
	mcp.AddTool(s, readOnlyTool("get_release_related_tags", "Get the related tags for a specific release."), toolHandler(client, tools.HandleGetReleaseRelatedTags))
	mcp.AddTool(s, readOnlyTool("get_release_tables", "Get the release tables for a specific release."), toolHandler(client, tools.HandleGetReleaseTables))
}

func registerSourceTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_sources", "Get all sources of economic data."), toolHandler(client, tools.HandleGetSources))
	mcp.AddTool(s, readOnlyTool("get_source", "Get a specific source of economic data."), toolHandler(client, tools.HandleGetSource))
	mcp.AddTool(s, readOnlyTool("get_source_releases", "Get the releases for a specific source."), toolHandler(client, tools.HandleGetSourceReleases))
}

func registerTagTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_tags", "Get FRED tags."), toolHandler(client, tools.HandleGetTags))
	mcp.AddTool(s, readOnlyTool("get_related_tags", "Get related FRED tags for one or more tag names."), toolHandler(client, tools.HandleGetRelatedTags))
	mcp.AddTool(s, readOnlyTool("get_tags_series", "Get the series matching one or more tags."), toolHandler(client, tools.HandleGetTagsSeries))
}

func registerSearchTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("search_series", "Search for economic data series by text."), toolHandler(client, tools.HandleSearchSeries))
	mcp.AddTool(s, readOnlyTool("get_release_series", "Get the series for a specific release."), toolHandler(client, tools.HandleGetReleaseSeries))
	mcp.AddTool(s, readOnlyTool("get_category_series", "Get the series for a specific category."), toolHandler(client, tools.HandleGetCategorySeries))
}

func registerGeoFREDTools(s *mcp.Server, client *fred.Client) {
	mcp.AddTool(s, readOnlyTool("get_series_group", "Get the series group for a GeoFRED series."), toolHandler(client, tools.HandleGetSeriesGroup))
	mcp.AddTool(s, readOnlyTool("get_series_data", "Get GeoFRED series map data."), toolHandler(client, tools.HandleGetSeriesData))
	mcp.AddTool(s, readOnlyTool("get_regional_data", "Get GeoFRED regional map data."), toolHandler(client, tools.HandleGetRegionalData))
}

func readOnlyTool(name, description string) *mcp.Tool {
	return &mcp.Tool{
		Name:        name,
		Description: description,
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func toolHandler[In any](client *fred.Client, fn toolFn[In]) func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		return fn(ctx, client, req, in)
	}
}
