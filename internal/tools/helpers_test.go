package tools_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shanehull/go-fred"
)

type toolCaller[In any] func(ctx context.Context, client *fred.Client, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error)

// call invokes a tool handler directly, mirroring how mcp.AddTool wraps
// handler errors and output.
func call[In any](ctx context.Context, client *fred.Client, fn toolCaller[In], in In) *mcp.CallToolResult {
	res, out, err := fn(ctx, client, nil, in)
	if err != nil {
		var errRes mcp.CallToolResult
		errRes.SetError(err)
		return &errRes
	}
	if res == nil {
		res = &mcp.CallToolResult{}
	}
	if out != nil {
		b, err := json.Marshal(out)
		if err != nil {
			var errRes mcp.CallToolResult
			errRes.SetError(err)
			return &errRes
		}
		res.Content = []mcp.Content{&mcp.TextContent{Text: string(b)}}
	}
	return res
}

func intPtr(v int) *int { return &v }

func newTestClient(t *testing.T, mock *httptest.Server) *fred.Client {
	t.Helper()
	opts := []fred.ClientOption{
		fred.WithAPIKey("test"),
		fred.WithHTTPClient(http.DefaultClient),
		fred.WithBaseURL("http://localhost"),
	}
	if mock != nil {
		opts[1] = fred.WithHTTPClient(mock.Client())
		opts[2] = fred.WithBaseURL(mock.URL)
	}
	client, err := fred.New(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func errorMock(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func assertTextContains(t *testing.T, result *mcp.CallToolResult, substr string) {
	t.Helper()
	if result.Content == nil {
		t.Fatalf("expected content containing %q, got nil", substr)
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(text.Text, substr) {
		t.Errorf("expected %q in result, got: %s", substr, text.Text)
	}
}

func assertIsError(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if !result.IsError {
		t.Fatal("expected error result, got success")
	}
}
