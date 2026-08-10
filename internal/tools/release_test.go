package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleGetRelease_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"releases":[{"id":53,"name":"Gross Domestic Product (GDP)","press_release":true}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetRelease, tools.ReleaseIDInput{ReleaseID: intPtr(53)})
	assertTextContains(t, result, "GDP")
}

func TestHandleGetRelease_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetRelease, tools.ReleaseIDInput{})
	assertIsError(t, result)
}

func TestHandleGetRelease_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetRelease, tools.ReleaseIDInput{ReleaseID: intPtr(999)})
	assertIsError(t, result)
}

func TestHandleGetReleases_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"releases":[{"id":53,"name":"GDP"},{"id":10,"name":"Industrial Production"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetReleases, tools.ReleaseListOptions{})
	assertTextContains(t, result, "GDP")
	assertTextContains(t, result, "Industrial")
}

func TestHandleGetReleaseDates_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "release_id=53") {
			t.Errorf("expected release_id=53, got %s", r.URL.String())
		}
		_, _ = w.Write([]byte(`{"release_dates":[{"release_id":53,"date":"2024-01-15"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetReleaseDates, tools.ReleaseIDDatesInput{ReleaseID: intPtr(53)})
	assertTextContains(t, result, "2024")
}

func TestHandleGetReleaseDates_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetReleaseDates, tools.ReleaseIDDatesInput{})
	assertIsError(t, result)
}

func TestHandleGetReleaseTables_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"elements":{"1":{"element_id":1,"name":"Table 1","level":"0"}}}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetReleaseTables, tools.ReleaseTablesInput{ReleaseID: intPtr(53)})
	assertTextContains(t, result, "Table 1")
}

func TestHandleGetReleaseTables_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetReleaseTables, tools.ReleaseTablesInput{})
	assertIsError(t, result)
}
