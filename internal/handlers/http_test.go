package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shanehull/fred-mcp/internal/config"
	"github.com/shanehull/fred-mcp/internal/handlers"
)

func testConfig() *config.Config {
	return &config.Config{
		OAuthIssuer:       "https://accounts.google.com",
		OAuthJwksURL:      "https://www.googleapis.com/oauth2/v3/certs",
		OAuthAuthorizeURL: "https://accounts.google.com/o/oauth2/v2/auth",
		OAuthTokenURL:     "https://oauth2.googleapis.com/token",
		OAuthAudience:     "test-client-id",
		OAuthClientSecret: "test-secret",
		OAuthCallbackURL:  "https://fred.example.com/oauth/callback",
	}
}

func mustQuery(t *testing.T, rawURL string) url.Values {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse url %q: %v", rawURL, err)
	}
	return parsed.Query()
}

func authorizeState(t *testing.T, cfg *config.Config, clientRedirect, clientState string) string {
	t.Helper()
	h := handlers.HandleAuthorizeProxy(cfg)

	target := "/authorize?client_id=x&redirect_uri=" + url.QueryEscape(clientRedirect) + "&state=" + url.QueryEscape(clientState)
	req := httptest.NewRequest("GET", target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302 from authorize, got %d", rec.Code)
	}
	return mustQuery(t, rec.Header().Get("Location")).Get("state")
}

func TestHandleDiscovery(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleDiscovery(cfg)

	req := httptest.NewRequest("GET", "/.well-known/oauth-protected-resource/mcp", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)

	if body["resource_name"] != "fred-mcp" {
		t.Errorf("expected resource_name=fred-mcp, got %v", body["resource_name"])
	}
	if body["client_id"] != "test-client-id" {
		t.Errorf("expected client_id=test-client-id, got %v", body["client_id"])
	}
	if _, ok := body["authorization_servers"]; !ok {
		t.Error("missing authorization_servers")
	}
}

func TestHandleDiscovery_RootPath(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleDiscovery(cfg)

	req := httptest.NewRequest("GET", "/.well-known/oauth-protected-resource", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)

	resource, _ := body["resource"].(string)
	if !strings.HasSuffix(resource, "/mcp") {
		t.Errorf("expected resource ending in /mcp, got %s", resource)
	}
}

func TestHandleAuthServerDiscovery(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleAuthServerDiscovery(cfg)

	req := httptest.NewRequest("GET", "/.well-known/oauth-authorization-server", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&body)

	if _, ok := body["registration_endpoint"]; !ok {
		t.Error("missing registration_endpoint")
	}
	if _, ok := body["authorization_endpoint"]; !ok {
		t.Error("missing authorization_endpoint")
	}
	if _, ok := body["token_endpoint"]; !ok {
		t.Error("missing token_endpoint")
	}
}

func TestHandleRegistration(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleRegistration(cfg)

	body := strings.NewReader(`{"redirect_uris":["http://localhost:3000"],"scope":"openid"}`)
	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}

	var res map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&res)

	if res["client_id"] != "test-client-id" {
		t.Errorf("expected client_id=test-client-id, got %v", res["client_id"])
	}
	if res["token_endpoint_auth_method"] != "none" {
		t.Errorf("expected token_endpoint_auth_method=none, got %v", res["token_endpoint_auth_method"])
	}
}

func TestHandleAuthorizeProxy_InjectScope(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleAuthorizeProxy(cfg)

	req := httptest.NewRequest("GET", "/authorize?client_id=x&redirect_uri=http://localhost", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, cfg.OAuthAuthorizeURL) {
		t.Errorf("expected redirect to %s, got %s", cfg.OAuthAuthorizeURL, loc)
	}
	if !strings.Contains(loc, "scope=openid+email+profile") {
		t.Errorf("expected scope injected, got %s", loc)
	}
}

func TestHandleAuthorizeProxy_RequestsOfflineAccess(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleAuthorizeProxy(cfg)

	req := httptest.NewRequest("GET", "/authorize?client_id=x&redirect_uri=http://localhost", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	query := mustQuery(t, rec.Header().Get("Location"))
	if got := query.Get("access_type"); got != "offline" {
		t.Errorf("expected access_type=offline, got %q", got)
	}
	if got := query.Get("prompt"); got != "consent" {
		t.Errorf("expected prompt=consent, got %q", got)
	}
}

func TestHandleAuthorizeProxy_PreservesOfflineAccessParams(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleAuthorizeProxy(cfg)

	req := httptest.NewRequest("GET", "/authorize?client_id=x&redirect_uri=http://localhost&access_type=online&prompt=none", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	query := mustQuery(t, rec.Header().Get("Location"))
	if got := query.Get("access_type"); got != "online" {
		t.Errorf("expected access_type preserved, got %q", got)
	}
	if got := query.Get("prompt"); got != "none" {
		t.Errorf("expected prompt preserved, got %q", got)
	}
}

func TestHandleAuthorizeProxy_RewritesRedirectURIAndState(t *testing.T) {
	cfg := testConfig()
	state := authorizeState(t, cfg, "http://localhost:3000/cb", "client-state")

	if state == "" || state == "client-state" {
		t.Errorf("expected wrapped state, got %q", state)
	}
}

func TestHandleAuthorizeProxy_RequiresRedirectURI(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleAuthorizeProxy(cfg)

	req := httptest.NewRequest("GET", "/authorize?client_id=x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleAuthorizeProxy_FallbackCallbackURL(t *testing.T) {
	cfg := testConfig()
	cfg.OAuthCallbackURL = ""
	h := handlers.HandleAuthorizeProxy(cfg)

	req := httptest.NewRequest("GET", "https://fred.example.com/authorize?client_id=x&redirect_uri="+url.QueryEscape("http://localhost/cb"), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	got := mustQuery(t, rec.Header().Get("Location")).Get("redirect_uri")
	if got != "https://fred.example.com/oauth/callback" {
		t.Errorf("expected callback derived from request host, got %s", got)
	}
}

func TestHandleAuthorizeCallback_RoundTrip(t *testing.T) {
	cfg := testConfig()
	clientRedirect := "http://localhost:3000/cb"
	state := authorizeState(t, cfg, clientRedirect, "client-state")

	h := handlers.HandleAuthorizeCallback(cfg)
	req := httptest.NewRequest("GET", "/oauth/callback?code=auth-code&state="+url.QueryEscape(state), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	loc, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse location: %v", err)
	}
	if got := loc.Scheme + "://" + loc.Host + loc.Path; got != clientRedirect {
		t.Errorf("expected redirect to %s, got %s", clientRedirect, got)
	}

	q := loc.Query()
	if q.Get("code") != "auth-code" {
		t.Errorf("expected code=auth-code, got %s", q.Get("code"))
	}
	if q.Get("state") != "client-state" {
		t.Errorf("expected restored state=client-state, got %s", q.Get("state"))
	}
	if _, ok := q["iss"]; ok {
		t.Error("iss must not be forwarded to the client")
	}
}

func TestHandleAuthorizeCallback_ErrorPassthrough(t *testing.T) {
	cfg := testConfig()
	state := authorizeState(t, cfg, "http://localhost:3000/cb", "client-state")

	h := handlers.HandleAuthorizeCallback(cfg)
	req := httptest.NewRequest("GET", "/oauth/callback?error=access_denied&error_description=Denied&state="+url.QueryEscape(state), nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	q := mustQuery(t, rec.Header().Get("Location"))
	if q.Get("error") != "access_denied" {
		t.Errorf("expected error=access_denied, got %s", q.Get("error"))
	}
	if q.Get("error_description") != "Denied" {
		t.Errorf("expected error_description=Denied, got %s", q.Get("error_description"))
	}
	if q.Get("state") != "client-state" {
		t.Errorf("expected restored state=client-state, got %s", q.Get("state"))
	}
	if _, ok := q["iss"]; ok {
		t.Error("iss must not be forwarded to the client")
	}
}

func TestHandleAuthorizeCallback_RejectsTamperedState(t *testing.T) {
	cfg := testConfig()
	state := authorizeState(t, cfg, "http://localhost:3000/cb", "client-state")

	h := handlers.HandleAuthorizeCallback(cfg)
	for name, raw := range map[string]string{
		"tampered":   state + "x",
		"malformed":  "not-a-token",
		"empty":      "",
		"no_payload": ".",
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/oauth/callback?code=auth-code&state="+url.QueryEscape(raw), nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rec.Code)
			}
		})
	}
}

func TestHandleTokenProxy_RewritesRedirectURIForAuthCode(t *testing.T) {
	var got url.Values
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"t"}`))
	}))
	defer upstream.Close()

	cfg := testConfig()
	cfg.OAuthTokenURL = upstream.URL
	h := handlers.HandleTokenProxy(cfg)

	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {"auth-code"},
		"redirect_uri": {"http://localhost:3000/cb"},
	}
	req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got.Get("redirect_uri") != cfg.OAuthCallbackURL {
		t.Errorf("expected redirect_uri=%s, got %s", cfg.OAuthCallbackURL, got.Get("redirect_uri"))
	}
}

func TestHandleTokenProxy_LeavesRefreshGrantUntouched(t *testing.T) {
	var got url.Values
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		got = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"t"}`))
	}))
	defer upstream.Close()

	cfg := testConfig()
	cfg.OAuthTokenURL = upstream.URL
	h := handlers.HandleTokenProxy(cfg)

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"refresh-token"},
	}
	req := httptest.NewRequest("POST", "/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got.Get("redirect_uri") != "" {
		t.Errorf("expected no redirect_uri for refresh grant, got %s", got.Get("redirect_uri"))
	}
}

func TestHandleTokenProxy_MethodNotAllowed(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleTokenProxy(cfg)

	req := httptest.NewRequest("GET", "/token", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestHandleConfig(t *testing.T) {
	cfg := testConfig()
	h := handlers.HandleConfig(cfg)

	req := httptest.NewRequest("GET", "/config", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)

	if body["type"] != "oauth" {
		t.Errorf("expected type=oauth, got %s", body["type"])
	}
	if body["clientId"] != "test-client-id" {
		t.Errorf("expected clientId=test-client-id, got %s", body["clientId"])
	}
}
