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
	result, _ := tools.HandleGetRelease(context.Background(), client, toolRequest("release_id", "53"))
	assertTextContains(t, result, "GDP")
}

func TestHandleGetRelease_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetRelease(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetRelease_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetRelease(context.Background(), client, toolRequest("release_id", "999"))
	assertIsError(t, result)
}

func TestHandleGetReleases_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"releases":[{"id":53,"name":"GDP"},{"id":10,"name":"Industrial Production"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetReleases(context.Background(), client, toolRequest())
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
	result, _ := tools.HandleGetReleaseDates(context.Background(), client, toolRequest("release_id", "53"))
	assertTextContains(t, result, "2024")
}

func TestHandleGetReleaseDates_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetReleaseDates(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetReleaseTables_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"elements":{"1":{"element_id":1,"name":"Table 1","level":"0"}}}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetReleaseTables(context.Background(), client, toolRequest("release_id", "53"))
	assertTextContains(t, result, "Table 1")
}

func TestHandleGetReleaseTables_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetReleaseTables(context.Background(), client, toolRequest())
	assertIsError(t, result)
}
