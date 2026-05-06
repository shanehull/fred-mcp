package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/fred-mcp/internal/tools"
	"github.com/shanehull/go-fred"
)

func TestHandleGetSeriesInfo_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "series_id=GDP") {
			t.Errorf("expected series_id=GDP in URL, got %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product","frequency":"Quarterly"}]}`))
	}))
	defer mock.Close()

	client, err := fred.New(
		fred.WithAPIKey("test"),
		fred.WithHTTPClient(mock.Client()),
		fred.WithBaseURL(mock.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"series_id": "GDP",
			},
		},
	}

	result, err := tools.HandleGetSeriesInfo(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil {
		t.Fatal("expected content in result")
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("expected text content")
	}
	if !strings.Contains(text.Text, "Gross Domestic Product") {
		t.Errorf("expected GDP in result, got %s", text.Text)
	}
}

func TestHandleGetSeriesInfo_MissingParam(t *testing.T) {
	client, err := fred.New(
		fred.WithAPIKey("test"),
		fred.WithHTTPClient(http.DefaultClient),
		fred.WithBaseURL("http://localhost"),
	)
	if err != nil {
		t.Fatal(err)
	}

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := tools.HandleGetSeriesInfo(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error result")
	}
}

func TestHandleGetSeriesInfo_ApiError(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":400,"error_message":"Bad request: invalid series_id"}`))
	}))
	defer mock.Close()

	client, err := fred.New(
		fred.WithAPIKey("test"),
		fred.WithHTTPClient(mock.Client()),
		fred.WithBaseURL(mock.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"series_id": "INVALID",
			},
		},
	}

	result, err := tools.HandleGetSeriesInfo(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error result")
	}
}

func TestHandleGetSeriesObservations_WithOptions(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		if !strings.Contains(url, "series_id=DGS20") {
			t.Errorf("expected series_id=DGS20 in URL, got %s", url)
		}
		if !strings.Contains(url, "units=lin") {
			t.Errorf("expected units=lin, got %s", url)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"observations":[{"date":"2024-01-01","value":"3.5"}]}`))
	}))
	defer mock.Close()

	client, err := fred.New(
		fred.WithAPIKey("test"),
		fred.WithHTTPClient(mock.Client()),
		fred.WithBaseURL(mock.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"series_id": "DGS20",
				"units":     "lin",
			},
		},
	}

	result, err := tools.HandleGetSeriesObservations(context.Background(), client, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content == nil {
		t.Fatal("expected content")
	}
}
