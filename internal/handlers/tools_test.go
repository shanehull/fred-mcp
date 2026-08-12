package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/fred-mcp/internal/handlers"
	"github.com/shanehull/go-fred"
)

func TestToolStructuredContentIsObject(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.String(), "observations") {
			_, _ = w.Write([]byte(`{"observations":[{"date":"2024-01-01","value":"3.5"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"seriess":[{"id":"UNRATE","title":"Unemployment Rate"}]}`))
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

	s := mcp.NewServer(&mcp.Implementation{Name: "fred-mcp", Version: "test"}, nil)
	handlers.RegisterTools(s, client)

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := s.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatal(err)
	}
	session, err := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0"}, nil).Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		toolName string
		args     any
	}{
		{"list output", "get_series_observations", map[string]any{"series_id": "UNRATE"}},
		{"object output", "get_series_info", map[string]any{"series_id": "UNRATE"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tt.toolName, Arguments: tt.args})
			if err != nil {
				t.Fatal(err)
			}
			assertObjectStructuredContent(t, res)
		})
	}
}

func assertObjectStructuredContent(t *testing.T, res *mcp.CallToolResult) {
	t.Helper()
	if res.StructuredContent == nil {
		t.Fatal("expected structuredContent, got nil")
	}
	data, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structuredContent: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("expected structuredContent to be a JSON object, got: %s", data)
	}
}
