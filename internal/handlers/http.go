package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/shanehull/fred-mcp/internal/config"
)

// fallbackStateKey signs authorize state when no client secret is configured.
var fallbackStateKey = func() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(fmt.Sprintf("handlers: generate OAuth state key: %v", err))
	}
	return key
}()

type authorizeState struct {
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state,omitempty"`
}

func stateKey(cfg *config.Config) []byte {
	if cfg.OAuthClientSecret != "" {
		return []byte(cfg.OAuthClientSecret)
	}
	return fallbackStateKey
}

func signState(key []byte, payload authorizeState) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(data) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func verifyState(key []byte, token string) (authorizeState, error) {
	var payload authorizeState
	dataPart, sigPart, ok := strings.Cut(token, ".")
	if !ok {
		return payload, errors.New("malformed state")
	}
	data, err := base64.RawURLEncoding.DecodeString(dataPart)
	if err != nil {
		return payload, errors.New("malformed state payload")
	}
	sig, err := base64.RawURLEncoding.DecodeString(sigPart)
	if err != nil {
		return payload, errors.New("malformed state signature")
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return payload, errors.New("state signature mismatch")
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return payload, errors.New("malformed state payload")
	}
	if payload.RedirectURI == "" {
		return payload, errors.New("state missing redirect_uri")
	}
	return payload, nil
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func callbackURL(cfg *config.Config, r *http.Request) string {
	if cfg.OAuthCallbackURL != "" {
		return cfg.OAuthCallbackURL
	}
	return baseURL(r) + "/oauth/callback"
}

func HandleDiscovery(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		host := baseURL(r)

		resourcePath := strings.TrimPrefix(r.URL.Path, "/.well-known/oauth-protected-resource")
		if resourcePath == "" {
			resourcePath = "/mcp"
		}

		slog.Info("discovery request", "host", host, "path", r.URL.Path, "resource_path", resourcePath)

		res := map[string]interface{}{
			"resource":                 fmt.Sprintf("%s%s", strings.TrimSuffix(host, "/"), resourcePath),
			"resource_name":            "fred-mcp",
			"authorization_servers":    []string{host},
			"client_id":                cfg.OAuthAudience,
			"clientId":                 cfg.OAuthAudience,
			"bearer_methods_supported": []string{"header"},
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			slog.Error("failed to encode discovery response", "error", err)
		}
	}
}

func HandleAuthServerDiscovery(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		host := baseURL(r)

		res := map[string]interface{}{
			"issuer":                           host,
			"authorization_endpoint":           fmt.Sprintf("%s/authorize", host),
			"token_endpoint":                   fmt.Sprintf("%s/token", host),
			"registration_endpoint":            fmt.Sprintf("%s/register", host),
			"jwks_uri":                         cfg.OAuthJwksURL,
			"response_types_supported":         []string{"code"},
			"grant_types_supported":            []string{"authorization_code", "refresh_token"},
			"scopes_supported":                 []string{"openid", "email", "profile"},
			"code_challenge_methods_supported": []string{"S256"},
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			slog.Error("failed to encode auth server discovery response", "error", err)
		}
	}
}

func HandleRegistration(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req map[string]interface{}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&req)
		}

		redirectURIs, _ := req["redirect_uris"].([]interface{})
		if redirectURIs == nil {
			redirectURIs = []interface{}{"http://localhost"}
		}

		scope, _ := req["scope"].(string)
		if scope == "" {
			scope = "openid email profile"
		}

		res := map[string]interface{}{
			"client_id":                  cfg.OAuthAudience,
			"client_id_issued_at":        0,
			"client_secret_expires_at":   0,
			"redirect_uris":              redirectURIs,
			"scope":                      scope,
			"token_endpoint_auth_method": "none",
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			slog.Error("failed to encode registration response", "error", err)
		}
	}
}

// redirectAllowed reports whether an authorization code may be delivered to raw.
// MCP clients use loopback redirects; every other URI must be listed in
// OAUTH_REDIRECT_URIS.
func redirectAllowed(cfg *config.Config, raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.Fragment != "" {
		return false
	}
	for _, allowed := range cfg.OAuthRedirectURIs {
		if raw == allowed {
			return true
		}
	}
	if u.Scheme != "http" {
		return false
	}
	switch u.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	return false
}

func HandleAuthorizeProxy(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("scope") == "" {
			params.Set("scope", "openid email profile")
		}
		// Google issues a refresh token only for offline access, and only repeats
		// it when consent is explicit. Without both, a login cannot outlive its
		// access token.
		if params.Get("access_type") == "" {
			params.Set("access_type", "offline")
		}
		if params.Get("prompt") == "" {
			params.Set("prompt", "consent")
		}

		clientRedirect := params.Get("redirect_uri")
		if clientRedirect == "" {
			http.Error(w, "Missing redirect_uri", http.StatusBadRequest)
			return
		}
		if !redirectAllowed(cfg, clientRedirect) {
			http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
			return
		}

		state, err := signState(stateKey(cfg), authorizeState{
			RedirectURI: clientRedirect,
			State:       params.Get("state"),
		})
		if err != nil {
			slog.Error("authorize proxy: failed to sign state", "error", err)
			http.Error(w, "Failed to encode state", http.StatusInternalServerError)
			return
		}

		params.Set("redirect_uri", callbackURL(cfg, r))
		params.Set("state", state)

		http.Redirect(w, r, cfg.OAuthAuthorizeURL+"?"+params.Encode(), http.StatusFound)
	}
}

// HandleAuthorizeCallback terminates the identity provider redirect, restores the
// client's redirect_uri and state, and returns the response without forwarding iss.
func HandleAuthorizeCallback(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		payload, err := verifyState(stateKey(cfg), params.Get("state"))
		if err != nil {
			slog.Warn("oauth callback: invalid state", "error", err)
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		redirect, err := url.Parse(payload.RedirectURI)
		if err != nil {
			slog.Warn("oauth callback: invalid redirect_uri in state", "error", err)
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		out := redirect.Query()
		if payload.State != "" {
			out.Set("state", payload.State)
		}
		if code := params.Get("code"); code != "" {
			out.Set("code", code)
		}
		if errCode := params.Get("error"); errCode != "" {
			out.Set("error", errCode)
			if desc := params.Get("error_description"); desc != "" {
				out.Set("error_description", desc)
			}
		}
		redirect.RawQuery = out.Encode()

		http.Redirect(w, r, redirect.String(), http.StatusFound)
	}
}

func HandleTokenProxy(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		form := url.Values{}
		for k, v := range r.Form {
			form[k] = v
		}
		form.Set("client_id", cfg.OAuthAudience)
		if cfg.OAuthClientSecret != "" {
			form.Set("client_secret", cfg.OAuthClientSecret)
		}
		if form.Get("grant_type") == "authorization_code" {
			form.Set("redirect_uri", callbackURL(cfg, r))
		}

		resp, err := http.PostForm(cfg.OAuthTokenURL, form)
		if err != nil {
			slog.Error("token proxy: upstream request failed", "error", err)
			http.Error(w, "Token exchange failed", http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			slog.Error("token proxy: failed to write response", "error", err)
		}
	}
}

func HandleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}

func HandleConfig(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := map[string]string{
			"type":     "oauth",
			"issuer":   cfg.OAuthIssuer,
			"clientId": cfg.OAuthAudience,
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			slog.Error("failed to encode config response", "error", err)
		}
	}
}
