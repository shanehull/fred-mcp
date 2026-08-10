package tools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSeriesInfo(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesIDInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	result, err := client.GetSeriesInfo(ctx, in.SeriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesObservations(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesObservationInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildObservationOptions(in.ObservationOptions)
	result, err := client.GetSeriesObservations(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesAllReleases(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesObservationInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildObservationOptions(in.ObservationOptions)
	result, err := client.GetSeriesAllReleases(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesFirstRelease(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesObservationInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildObservationOptions(in.ObservationOptions)
	result, err := client.GetSeriesFirstRelease(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesAsOf(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesAsOfInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	if in.AsOfDate == "" {
		return nil, nil, errors.New("as_of_date is required")
	}
	asOf, err := time.Parse("2006-01-02", in.AsOfDate)
	if err != nil {
		return nil, nil, errors.New("invalid as_of_date format, expected YYYY-MM-DD")
	}
	opts := buildObservationOptions(in.ObservationOptions)
	result, err := client.GetSeriesAsOf(ctx, in.SeriesID, asOf, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesVintageDates(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesObservationInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildObservationOptions(in.ObservationOptions)
	result, err := client.GetSeriesVintageDates(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesCategories(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesIDInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	result, err := client.GetSeriesCategories(ctx, in.SeriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesRelease(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesIDInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	result, err := client.GetSeriesRelease(ctx, in.SeriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesIDTagInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetSeriesTags(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleSearchSeriesTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SearchSeriesTagsInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesSearchText == "" {
		return nil, nil, errors.New("series_search_text is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.SearchSeriesTags(ctx, in.SeriesSearchText, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleSearchSeriesRelatedTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SearchSeriesRelatedTagsInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesSearchText == "" {
		return nil, nil, errors.New("series_search_text is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.SearchSeriesRelatedTags(ctx, in.SeriesSearchText, parseStringList(in.TagNames), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesUpdates(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in UpdateOptions) (*mcp.CallToolResult, any, error) {
	opts := buildUpdateOptions(in)
	result, err := client.GetSeriesUpdates(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
