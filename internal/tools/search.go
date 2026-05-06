package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleSearchSeries(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := req.RequireString("search_text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildSearchOptions(req)
	result, err := client.SearchSeries(ctx, text, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseSeries(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildSearchOptions(req)
	result, err := client.GetReleaseSeries(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetCategorySeries(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildSearchOptions(req)
	result, err := client.GetCategorySeries(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
