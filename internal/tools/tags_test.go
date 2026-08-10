package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleGetTags_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tags":[{"name":"gdp","group_id":"geo","series_count":500},{"name":"inflation","group_id":"other","series_count":200}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetTags, tools.TagsInput{})
	assertTextContains(t, result, "gdp")
	assertTextContains(t, result, "inflation")
}

func TestHandleGetRelatedTags_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tags":[{"name":"business","group_id":" ","series_count":300}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetRelatedTags, tools.RelatedTagsInput{TagNames: "gdp"})
	assertTextContains(t, result, "business")
}

func TestHandleGetRelatedTags_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetRelatedTags, tools.RelatedTagsInput{})
	assertIsError(t, result)
}

func TestHandleGetTagsSeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetTagsSeries, tools.TagsSeriesInput{TagNames: "gdp"})
	assertTextContains(t, result, "GDP")
}

func TestHandleGetTagsSeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetTagsSeries, tools.TagsSeriesInput{})
	assertIsError(t, result)
}
