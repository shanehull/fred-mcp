package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSources(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildSourceOptions(req)
	result, err := client.GetSources(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSource(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "source_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetSource(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetSourceReleases(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "source_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildReleaseListOptions(req)
	result, err := client.GetSourceReleases(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
