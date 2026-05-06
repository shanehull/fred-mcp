package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetReleases(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildReleaseListOptions(req)
	result, err := client.GetReleases(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleasesDates(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildReleaseDateOptions(req)
	result, err := client.GetReleasesDates(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetRelease(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetRelease(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseDates(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildReleaseDateOptions(req)
	result, err := client.GetReleaseDates(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseSources(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetReleaseSources(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildTagOptions(req)
	result, err := client.GetReleaseTags(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseRelatedTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tagNames := parseStringList(req, "tag_names")
	opts := buildTagOptions(req)
	result, err := client.GetReleaseRelatedTags(ctx, id, tagNames, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetReleaseTables(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "release_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildTableOptions(req)
	result, err := client.GetReleaseTables(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
