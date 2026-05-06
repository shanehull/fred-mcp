package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetCategory(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetCategory(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetCategoryChildren(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetCategoryChildren(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetCategoryRelated(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	result, err := client.GetCategoryRelated(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetCategoryTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	opts := buildTagOptions(req)
	result, err := client.GetCategoryTags(ctx, id, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}

func HandleGetCategoryRelatedTags(ctx context.Context, client *fred.Client, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := GetInt(req, "category_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	tagNames := parseStringList(req, "tag_names")
	opts := buildTagOptions(req)
	result, err := client.GetCategoryRelatedTags(ctx, id, tagNames, opts...)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("FRED API error: %v", err)), nil
	}
	return MarshalResult(result)
}
