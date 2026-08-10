package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

func HandleGetReleases(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseListOptions) (*mcp.CallToolResult, any, error) {
	opts := buildReleaseListOptions(in)
	result, err := client.GetReleases(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleasesDates(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseDateOptions) (*mcp.CallToolResult, any, error) {
	opts := buildReleaseDateOptions(in)
	result, err := client.GetReleasesDates(ctx, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetRelease(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseIDInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	result, err := client.GetRelease(ctx, *in.ReleaseID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseDates(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseIDDatesInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	opts := buildReleaseDateOptions(in.ReleaseDateOptions)
	result, err := client.GetReleaseDates(ctx, *in.ReleaseID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseSources(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseIDInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	result, err := client.GetReleaseSources(ctx, *in.ReleaseID)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseIDTagInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetReleaseTags(ctx, *in.ReleaseID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseRelatedTags(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseRelatedTagsInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	opts := buildTagOptions(in.TagOptions)
	result, err := client.GetReleaseRelatedTags(ctx, *in.ReleaseID, parseStringList(in.TagNames), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}

func HandleGetReleaseTables(ctx context.Context, client *fred.Client, _ *mcp.CallToolRequest, in ReleaseTablesInput) (*mcp.CallToolResult, any, error) {
	if in.ReleaseID == nil {
		return nil, nil, errors.New("release_id is required")
	}
	opts := buildTableOptions(in.TableOptions)
	result, err := client.GetReleaseTables(ctx, *in.ReleaseID, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("FRED API error: %v", err)
	}
	return nil, result, nil
}
