package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in TagsInput) (*mcp.CallToolResult, any, error) {
	opts := buildTagOptions(in.TagOptions)
	if in.TagNames != "" {
		opts = append(opts, fred.WithTagSetNames(parseStringList(in.TagNames)...))
	}
	result, err := client.GetTags(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetRelatedTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in RelatedTagsInput) (*mcp.CallToolResult, any, error) {
	if in.TagNames == "" {
		return nil, nil, errors.New("tag_names is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetRelatedTags(ctx, parseStringList(in.TagNames), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetTagsSeries(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in TagsSeriesInput) (*mcp.CallToolResult, any, error) {
	if in.TagNames == "" {
		return nil, nil, errors.New("tag_names is required")
	}
	opts := buildSearchOptions(in.SearchOptions)
	if in.TagNames != "" {
		opts = append(opts, fred.WithTagNames(parseStringList(in.TagNames)...))
	}
	result, err := client.GetTagsSeries(ctx, parseStringList(in.TagNames), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
