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
	result := call(context.Background(), client, tools.HandleSearchSeries, tools.SearchSeriesInput{SearchText: "GDP"})
	assertTextContains(t, result, "Gross Domestic Product")
}

func TestHandleSearchSeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleSearchSeries, tools.SearchSeriesInput{})
	assertIsError(t, result)
}

func TestHandleSearchSeries_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleSearchSeries, tools.SearchSeriesInput{SearchText: ""})
	assertIsError(t, result)
}

func TestHandleGetReleaseSeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"GDP","title":"Gross Domestic Product"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetReleaseSeries, tools.ReleaseSeriesInput{ReleaseID: intPtr(53)})
	assertTextContains(t, result, "GDP")
}

func TestHandleGetReleaseSeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetReleaseSeries, tools.ReleaseSeriesInput{})
	assertIsError(t, result)
}

func TestHandleGetCategorySeries_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"seriess":[{"id":"CPI","title":"Consumer Price Index"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetCategorySeries, tools.CategorySeriesInput{CategoryID: intPtr(1)})
	assertTextContains(t, result, "Consumer Price Index")
}

func TestHandleGetCategorySeries_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetCategorySeries, tools.CategorySeriesInput{})
	assertIsError(t, result)
}
