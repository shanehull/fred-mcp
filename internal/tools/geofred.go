package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSeriesGroup(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesGroupInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	result, err := client.GetSeriesGroup(ctx, in.SeriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSeriesData(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SeriesDataInput) (*mcp.CallToolResult, any, error) {
	if in.SeriesID == "" {
		return nil, nil, errors.New("series_id is required")
	}
	opts := buildMapDataOptions(in.MapDataOptions)
	result, err := client.GetSeriesData(ctx, in.SeriesID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetRegionalData(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in RegionalDataInput) (*mcp.CallToolResult, any, error) {
	opts := buildRegionalDataOptions(in)
	result, err := client.GetRegionalData(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
