package tools_test

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/fred-mcp/internal/tools"
	"github.com/shanehull/go-fred"
)

func skipIfNoKey(t *testing.T) {
	t.Helper()
	if os.Getenv("FRED_API_KEY") == "" {
		t.Skip("FRED_API_KEY not set")
	}
}

func newLiveClient(t *testing.T) *fred.Client {
	t.Helper()
	client, err := fred.New()
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Series
func TestLive_GetSeriesInfo(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_series_info", "GDP") }
func TestLive_GetSeriesObservations(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetSeriesObservations(context.Background(), client, toolRequest(
		"series_id", "UNRATE",
		"limit", "3",
		"sort_order", "desc",
	))
	assertTextContainsLive(t, result)
}

// Category
func TestLive_GetCategory(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_category", "1") }
func TestLive_GetCategoryChildren(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetCategoryChildren(context.Background(), client, toolRequest("category_id", "0"))
	assertTextContainsLive(t, result)
}

// Release
func TestLive_GetRelease(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_release", "53") }
func TestLive_GetReleases(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetReleases(context.Background(), client, toolRequest("limit", "5"))
	assertTextContainsLive(t, result)
}

// Source
func TestLive_GetSource(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_source", "1") }
func TestLive_GetSources(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetSources(context.Background(), client, toolRequest("limit", "3"))
	assertTextContainsLive(t, result)
}

// Tags
func TestLive_GetTags(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetTags(context.Background(), client, toolRequest("limit", "3", "tag_names", "gdp,inflation"))
	assertTextContainsLive(t, result)
}

// Search
func TestLive_SearchSeries(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleSearchSeries(context.Background(), client, toolRequest("search_text", "GDP", "limit", "3"))
	assertTextContainsLive(t, result)
}

// GeoFRED
func TestLive_GetSeriesGroup(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result, _ := tools.HandleGetSeriesGroup(context.Background(), client, toolRequest("series_id", "WIPCPI"))
	assertTextContainsLive(t, result)
}

func testLive(t *testing.T, toolName string, id string) {
	t.Helper()
	client := newLiveClient(t)

	var result *mcp.CallToolResult
	var goErr error

	switch toolName {
	case "fred_get_series_info":
		result, goErr = tools.HandleGetSeriesInfo(context.Background(), client, toolRequest("series_id", id))
	case "fred_get_category":
		result, goErr = tools.HandleGetCategory(context.Background(), client, toolRequest("category_id", id))
	case "fred_get_release":
		result, goErr = tools.HandleGetRelease(context.Background(), client, toolRequest("release_id", id))
	case "fred_get_source":
		result, goErr = tools.HandleGetSource(context.Background(), client, toolRequest("source_id", id))
	default:
		t.Fatalf("unknown tool: %s", toolName)
	}
	if goErr != nil {
		t.Fatal(goErr)
	}
	assertTextContainsLive(t, result)
}

func assertTextContainsLive(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result.IsError {
		text := ""
		if len(result.Content) > 0 {
			if tc, ok := result.Content[0].(mcp.TextContent); ok {
				text = tc.Text
			}
		}
		t.Fatalf("API returned error: %s", text)
	}
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}
}
