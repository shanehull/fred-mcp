package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSeriesInfo(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetSeriesInfo(ctx, seriesID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesObservations(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildObservationOptions(req)
	result, err := client.GetSeriesObservations(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesAllReleases(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildObservationOptions(req)
	result, err := client.GetSeriesAllReleases(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesFirstRelease(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildObservationOptions(req)
	result, err := client.GetSeriesFirstRelease(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesAsOf(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	asOfStr, err := req.RequireString("as_of_date")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	asOf, err := time.Parse("2006-01-02", asOfStr)
	if err != nil {
		return mcp.NewToolResultError("invalid as_of_date format, expected YYYY-MM-DD"), nil
	}
	opts := buildObservationOptions(req)
	result, err := client.GetSeriesAsOf(ctx, seriesID, asOf, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesVintageDates(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildObservationOptions(req)
	result, err := client.GetSeriesVintageDates(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesCategories(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetSeriesCategories(ctx, seriesID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesRelease(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetSeriesRelease(ctx, seriesID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildTagOptions(req)
	result, err := client.GetSeriesTags(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleSearchSeriesTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := req.RequireString("series_search_text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildTagOptions(req)
	result, err := client.SearchSeriesTags(ctx, text, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleSearchSeriesRelatedTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := req.RequireString("series_search_text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tagNames := parseStringList(req, "tag_names")
	opts := buildTagOptions(req)
	result, err := client.SearchSeriesRelatedTags(ctx, text, tagNames, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesUpdates(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildUpdateOptions(req)
	result, err := client.GetSeriesUpdates(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
