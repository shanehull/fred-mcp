package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleSearchSeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product"},{"id":"GDPC1","title":"Real Gross Domestic Product"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleSearchSeries(context.Background(), client, toolRequest("search_text", "GDP"))
	assertTextContains(t, result, "Gross Domestic Product")
}

func TestHandleSearchSeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleSearchSeries(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleSearchSeries_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleSearchSeries(context.Background(), client, toolRequest("search_text", ""))
	assertIsError(t, result)
}

func TestHandleGetReleaseSeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetReleaseSeries(context.Background(), client, toolRequest("release_id", "53"))
	assertTextContains(t, result, "GDP")
}

func TestHandleGetReleaseSeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetReleaseSeries(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetCategorySeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"CPI","title":"Consumer Price Index"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetCategorySeries(context.Background(), client, toolRequest("category_id", "1"))
	assertTextContains(t, result, "Consumer Price Index")
}

func TestHandleGetCategorySeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetCategorySeries(context.Background(), client, toolRequest())
	assertIsError(t, result)
}
