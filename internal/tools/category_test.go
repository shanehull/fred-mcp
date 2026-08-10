package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shanehull/fred-mcp/internal/tools"
)

func TestHandleGetCategory_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "category_id=1") {
			t.Errorf("expected category_id=1, got %s", r.URL.String())
		}
		_, _ = w.Write([]byte(`{"categories":[{"id":1,"name":"Production & Business Activity","parent_id":0}]}`))
	}))
	defer mock.Close()

	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetCategory, tools.CategoryIDInput{CategoryID: intPtr(1)})
	assertTextContains(t, result, "Production")
}

func TestHandleGetCategory_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetCategory, tools.CategoryIDInput{})
	assertIsError(t, result)
}

func TestHandleGetCategory_ApiError(t *testing.T) {
	mock := errorMock(t, http.StatusBadRequest, `{"error_code":400,"error_message":"invalid"}`)
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetCategory, tools.CategoryIDInput{CategoryID: intPtr(999)})
	assertIsError(t, result)
}

func TestHandleGetCategoryChildren_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"categories":[{"id":2,"name":"Child Category","parent_id":1}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetCategoryChildren, tools.CategoryIDInput{CategoryID: intPtr(1)})
	assertTextContains(t, result, "Child")
}

func TestHandleGetCategoryChildren_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetCategoryChildren, tools.CategoryIDInput{})
	assertIsError(t, result)
}

func TestHandleGetCategoryTags_Success(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tags":[{"name":"gdp","group_id":"geo","series_count":100}]}`))
	}))
	defer mock.Close()
	client := newTestClient(t, mock)
	result := call(context.Background(), client, tools.HandleGetCategoryTags, tools.CategoryIDTagInput{CategoryID: intPtr(1)})
	assertTextContains(t, result, "gdp")
}

func TestHandleGetCategoryTags_MissingParam(t *testing.T) {
	client := newTestClient(t, nil)
	result := call(context.Background(), client, tools.HandleGetCategoryTags, tools.CategoryIDTagInput{})
	assertIsError(t, result)
}
