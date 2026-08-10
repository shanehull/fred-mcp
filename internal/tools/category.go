package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetCategory(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategoryIDInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	result, err := client.GetCategory(ctx, *in.CategoryID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetCategoryChildren(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategoryIDInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	result, err := client.GetCategoryChildren(ctx, *in.CategoryID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetCategoryRelated(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategoryIDInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	result, err := client.GetCategoryRelated(ctx, *in.CategoryID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetCategoryTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategoryIDTagInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetCategoryTags(ctx, *in.CategoryID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetCategoryRelatedTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategoryRelatedTagsInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetCategoryRelatedTags(ctx, *in.CategoryID, parseStringList(in.TagNames), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
