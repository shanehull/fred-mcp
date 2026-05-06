package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shanehull/fred-mcp/internal/config"
	"github.com/shanehull/fred-mcp/internal/middleware"
)

func newTestJWKS(t *testing.T) (*keyfunc.JWKS, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	givenKey := keyfunc.NewGivenRSA(&key.PublicKey, keyfunc.GivenKeyOptions{Algorithm: "RS256"})
	jwks := keyfunc.NewGiven(map[string]keyfunc.GivenKey{"test": givenKey})
	return jwks, key
}

func signJWT(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test"
	s, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func requestWithToken(token string) *http.Request {
	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuth_NilJWKS_PassThrough(t *testing.T) {
	cfg := &config.Config{OAuthAudience: "test"}
	handler := middleware.Auth(cfg, nil)(okHandler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/sse", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuth_MissingBearer(t *testing.T) {
	cfg := &config.Config{OAuthAudience: "test"}
	jwks, _ := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/sse", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
	if rec.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header")
	}
}

func TestAuth_ValidJWT(t *testing.T) {
	cfg := &config.Config{OAuthAudience: "test-client"}
	jwks, key := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	token := signJWT(t, key, jwt.MapClaims{"aud": "test-client"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, requestWithToken(token))

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuth_InvalidToken_Returns401(t *testing.T) {
	cfg := &config.Config{OAuthAudience: "test-client"}
	jwks, _ := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, requestWithToken("not-a-valid-token-at-all"))

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAuth_ValidJWT_AudienceMismatch(t *testing.T) {
	cfg := &config.Config{OAuthAudience: "correct-client"}
	jwks, key := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	token := signJWT(t, key, jwt.MapClaims{"aud": "wrong-client"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, requestWithToken(token))

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

func TestAuth_ValidJWT_EmailMatch(t *testing.T) {
	cfg := &config.Config{
		OAuthAudience: "test-client",
		AllowedEmail:  "alice@example.com",
	}
	jwks, key := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	token := signJWT(t, key, jwt.MapClaims{
		"aud":   "test-client",
		"email": "alice@example.com",
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, requestWithToken(token))

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuth_ValidJWT_EmailMismatch(t *testing.T) {
	cfg := &config.Config{
		OAuthAudience: "test-client",
		AllowedEmail:  "alice@example.com",
	}
	jwks, key := newTestJWKS(t)
	handler := middleware.Auth(cfg, jwks)(okHandler())

	token := signJWT(t, key, jwt.MapClaims{
		"aud":   "test-client",
		"email": "bob@example.com",
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, requestWithToken(token))

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}
