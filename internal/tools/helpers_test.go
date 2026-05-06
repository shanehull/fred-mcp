package tools_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shanehull/go-fred"
)

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

func toolRequest(kv ...string) mcp.CallToolRequest {
	args := make(map[string]interface{})
	for i := 0; i < len(kv)-1; i += 2 {
		args[kv[i]] = kv[i+1]
	}
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: args,
		},
	}
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
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
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
