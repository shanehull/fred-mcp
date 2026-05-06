package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleGetSource_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sources":[{"id":1,"name":"Board of Governors of the Federal Reserve System","link":"http://www.federalreserve.gov/"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSource(context.Background(), client, toolRequest("source_id", "1"))
	assertTextContains(t, result, "Federal Reserve")
}

func TestHandleGetSource_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetSource(context.Background(), client, toolRequest())
	assertIsError(t, result)
}

func TestHandleGetSource_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSource(context.Background(), client, toolRequest("source_id", "999"))
	assertIsError(t, result)
}

func TestHandleGetSources_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sources":[{"id":1,"name":"Federal Reserve"},{"id":2,"name":"Bureau of Labor Statistics"}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSources(context.Background(), client, toolRequest())
	assertTextContains(t, result, "Labor")
}

func TestHandleGetSourceReleases_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"releases":[{"id":53,"name":"GDP","press_release":true}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result, _ := tools.HandleGetSourceReleases(context.Background(), client, toolRequest("source_id", "1"))
	assertTextContains(t, result, "GDP")
}

func TestHandleGetSourceReleases_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result, _ := tools.HandleGetSourceReleases(context.Background(), client, toolRequest())
	assertIsError(t, result)
}
