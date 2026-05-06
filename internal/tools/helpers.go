package tools

import (
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

func MarshalResult(v interface{}) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal result: " + err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
