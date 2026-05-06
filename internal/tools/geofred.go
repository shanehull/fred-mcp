package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSeriesGroup(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetSeriesGroup(ctx, seriesID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSeriesData(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seriesID, err := req.RequireString("series_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildMapDataOptions(req)
	result, err := client.GetSeriesData(ctx, seriesID, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetRegionalData(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildRegionalDataOptions(req)
	result, err := client.GetRegionalData(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
