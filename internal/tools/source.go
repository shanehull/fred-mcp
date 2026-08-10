package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetSources(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SourceOptions) (*mcp.CallToolResult, any, error) {
	opts := buildSourceOptions(in)
	result, err := client.GetSources(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSource(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SourceIDInput) (*mcp.CallToolResult, any, error) {
	if in.SourceID == nil {
		return nil, nil, errors.New("source_id is required")
	}
	result, err := client.GetSource(ctx, *in.SourceID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetSourceReleases(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in SourceIDReleasesInput) (*mcp.CallToolResult, any, error) {
	if in.SourceID == nil {
		return nil, nil, errors.New("source_id is required")
	}
	opts := buildReleaseListOptions(in.ReleaseListOptions)
	result, err := client.GetSourceReleases(ctx, *in.SourceID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
