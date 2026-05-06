package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleGetSeriesGroup_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"series_group":{"title":"Per Capita Personal Income","region_type":"state","series_group":"882","season":"NSA","units":"Dollars","frequency":"Annual"}}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSeriesGroup(context.Background(), client, toolRequest("series_id", "WIPCPI"))
	assertTextContains(t, result, "Personal Income")
}

func TestHandleGetSeriesGroup_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetSeriesGroup(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetSeriesGroup_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSeriesGroup(context.Background(), client, toolRequest("series_id", "INVALID"))
	assertIsError(t, result)
}

func TestHandleGetSeriesData_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"meta":{"title":"Per Capita Personal Income","region":"state","seasonality":"NSA","units":"Dollars","frequency":"Annual","date":"2012-01-01","data":{"2012-01-01":[{"region":"California","code":"CA","value":45000}]}}}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSeriesData(context.Background(), client, toolRequest("series_id", "WIPCPI"))
	assertTextContains(t, result, "California")
}

func TestHandleGetSeriesData_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetSeriesData(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetRegionalData_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"meta":{"title":"Per Capita Personal Income","region":"state","seasonality":"NSA","units":"Dollars","frequency":"Annual","date":"2013-01-01","data":{"2013-01-01":[{"region":"New York","code":"NY","value":55000}]}}}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetRegionalData(context.Background(), client, toolRequest("series_group", "882"))
	assertTextContains(t, result, "New York")
}
