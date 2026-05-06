package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts := buildTagOptions(req)
	result, err := client.GetTags(ctx, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetRelatedTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tagNames, err := req.RequireString("tag_names")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	names := parseStringListFromRaw(tagNames)
	opts := buildTagOptions(req)
	result, err := client.GetRelatedTags(ctx, names, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetTagsSeries(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tagNames, err := req.RequireString("tag_names")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	names := parseStringListFromRaw(tagNames)
	opts := buildSearchOptions(req)
	result, err := client.GetTagsSeries(ctx, names, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func parseStringListFromRaw(raw string) []string {
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
