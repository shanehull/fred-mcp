package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleSearchSeries(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SearchSeriesInput) (*mcp.CallToolResult, any, error) {
	if in.SearchText == "" {
		return nil, nil, errors.New("search_text is required")
	}
	opts := buildSearchOptions(in.SearchOptions)
	if in.TagNames != "" {
		opts = append(opts, fred.WithTagNames(parseStringList(in.TagNames)...))
	}
	result, err := client.SearchSeries(ctx, in.SearchText, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseSeries(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseSeriesInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	opts := buildSearchOptions(in.SearchOptions)
	if in.TagNames != "" {
		opts = append(opts, fred.WithTagNames(parseStringList(in.TagNames)...))
	}
	result, err := client.GetReleaseSeries(ctx, *in.ReleaseID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetCategorySeries(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in CategorySeriesInput) (*mcp.CallToolResult, any, error) {
	if in.CategoryID == nil {
		return nil, nil, errors.New("category_id is required")
	}
	opts := buildSearchOptions(in.SearchOptions)
	if in.TagNames != "" {
		opts = append(opts, fred.WithTagNames(parseStringList(in.TagNames)...))
	}
	result, err := client.GetCategorySeries(ctx, *in.CategoryID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
