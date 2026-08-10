package tools_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	result := call(context.Background(), client, tools.HandleGetSeriesObservations, tools.SeriesObservationInput{
		SeriesID: "UNRATE",
		ObservationOptions: tools.ObservationOptions{
			Limit:     3,
			SortOrder: "desc",
		},
	})
	assertTextContainsLive(t, result)
}

// Category
func TestLive_GetCategory(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_category", "1") }
func TestLive_GetCategoryChildren(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleGetCategoryChildren, tools.CategoryIDInput{CategoryID: intPtr(0)})
	assertTextContainsLive(t, result)
}

// Release
func TestLive_GetRelease(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_release", "53") }
func TestLive_GetReleases(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleGetReleases, tools.ReleaseListOptions{Limit: 5})
	assertTextContainsLive(t, result)
}

// Source
func TestLive_GetSource(t *testing.T) { skipIfNoKey(t); testLive(t, "fred_get_source", "1") }
func TestLive_GetSources(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleGetSources, tools.SourceOptions{Limit: 3})
	assertTextContainsLive(t, result)
}

// Tags
func TestLive_GetTags(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleGetTags, tools.TagsInput{
		TagNames:   "gdp,inflation",
		TagOptions: tools.TagOptions{Limit: 3},
	})
	assertTextContainsLive(t, result)
}

// Search
func TestLive_SearchSeries(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleSearchSeries, tools.SearchSeriesInput{
		SearchText:   "GDP",
		SearchOptions: tools.SearchOptions{Limit: 3},
	})
	assertTextContainsLive(t, result)
}

// GeoFRED
func TestLive_GetSeriesGroup(t *testing.T) {
	skipIfNoKey(t)
	client := newLiveClient(t)
	result := call(context.Background(), client, tools.HandleGetSeriesGroup, tools.SeriesIDInput{SeriesID: "WIPCPI"})
	assertTextContainsLive(t, result)
}

func testLive(t *testing.T, toolName string, id string) {
	t.Helper()
	client := newLiveClient(t)

	var result *mcp.CallToolResult

	switch toolName {
	case "fred_get_series_info":
		result = call(context.Background(), client, tools.HandleGetSeriesInfo, tools.SeriesIDInput{SeriesID: id})
	case "fred_get_category":
		cid, _ := strconv.Atoi(id)
		result = call(context.Background(), client, tools.HandleGetCategory, tools.CategoryIDInput{CategoryID: &cid})
	case "fred_get_release":
		rid, _ := strconv.Atoi(id)
		result = call(context.Background(), client, tools.HandleGetRelease, tools.ReleaseIDInput{ReleaseID: &rid})
	case "fred_get_source":
		sid, _ := strconv.Atoi(id)
		result = call(context.Background(), client, tools.HandleGetSource, tools.SourceIDInput{SourceID: &sid})
	default:
		t.Fatalf("unknown tool: %s", toolName)
	}
	assertTextContainsLive(t, result)
}

func assertTextContainsLive(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result.IsError {
		text := ""
		if len(result.Content) > 0 {
			if tc, ok := result.Content[0].(*mcp.TextContent); ok {
				text = tc.Text
			}
		}
		t.Fatalf("API returned error: %s", text)
	}
	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}
}
