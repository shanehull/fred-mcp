package config_test

import (
	"os"
	"testing"

	"github.com/shanehull/fred-mcp/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("FRED_API_KEY")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("OAUTH_ISSUER")
	_ = os.Unsetenv("OAUTH_JWKS_URL")
	_ = os.Unsetenv("OAUTH_AUTHORIZE_URL")
	_ = os.Unsetenv("OAUTH_TOKEN_URL")
	_ = os.Unsetenv("PUBLIC_HOST")
	_ = os.Unsetenv("OAUTH_AUDIENCE")
	_ = os.Unsetenv("OAUTH_CLIENT_SECRET")
	_ = os.Unsetenv("OAUTH_ALLOWED_EMAILS")

	cfg := config.Load()

	if cfg.FredAPIKey != "" {
		t.Error("expected empty FredAPIKey by default")
	}
	if cfg.Port != "4000" {
		t.Errorf("expected Port=4000, got %s", cfg.Port)
	}
	if cfg.OAuthIssuer != "https://accounts.google.com" {
		t.Errorf("expected Google issuer, got %s", cfg.OAuthIssuer)
	}
	if cfg.OAuthJwksURL != "https://www.googleapis.com/oauth2/v3/certs" {
		t.Errorf("expected Google JWKS, got %s", cfg.OAuthJwksURL)
	}
	if cfg.OAuthAudience != "" {
		t.Errorf("expected empty OAuthAudience, got %s", cfg.OAuthAudience)
	}
}

func TestLoadOverrides(t *testing.T) {
	_ = os.Setenv("FRED_API_KEY", "test-key")
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("OAUTH_AUDIENCE", "test-client-id")
	_ = os.Setenv("OAUTH_CLIENT_SECRET", "test-secret")
	_ = os.Setenv("OAUTH_ALLOWED_EMAILS", "test@example.com,bob@example.com")
	_ = os.Setenv("PUBLIC_HOST", "https://example.com")
	defer func() { _ = os.Unsetenv("FRED_API_KEY") }()
	defer func() { _ = os.Unsetenv("PORT") }()
	defer func() { _ = os.Unsetenv("OAUTH_AUDIENCE") }()
	defer func() { _ = os.Unsetenv("OAUTH_CLIENT_SECRET") }()
	defer func() { _ = os.Unsetenv("OAUTH_ALLOWED_EMAILS") }()
	defer func() { _ = os.Unsetenv("PUBLIC_HOST") }()

	cfg := config.Load()

	if cfg.FredAPIKey != "test-key" {
		t.Errorf("expected FredAPIKey=test-key, got %s", cfg.FredAPIKey)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port=8080, got %s", cfg.Port)
	}
	if cfg.OAuthAudience != "test-client-id" {
		t.Errorf("expected OAuthAudience=test-client-id, got %s", cfg.OAuthAudience)
	}
	if cfg.OAuthClientSecret != "test-secret" {
		t.Errorf("expected OAuthClientSecret=test-secret, got %s", cfg.OAuthClientSecret)
	}
	if cfg.AllowedEmails == nil {
		t.Error("expected AllowedEmails to be set")
	} else {
		if len(cfg.AllowedEmails) != 2 {
			t.Errorf("expected 2 allowed emails, got %d", len(cfg.AllowedEmails))
		}
		if cfg.AllowedEmails[0] != "test@example.com" {
			t.Errorf("expected first email=test@example.com, got %s", cfg.AllowedEmails[0])
		}
	}
	if cfg.PublicHost != "https://example.com" {
		t.Errorf("expected PublicHost=https://example.com, got %s", cfg.PublicHost)
	}
}
