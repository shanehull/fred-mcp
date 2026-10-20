package config

import (
	"os"
	"strings"
)

type Config struct {
	FredAPIKey        string
	Port              string
	OAuthIssuer       string
	OAuthJwksURL      string
	OAuthAuthorizeURL string
	OAuthTokenURL     string
	OAuthAudience     string
	OAuthClientSecret string
	OAuthCallbackURL  string
	OAuthRedirectURIs []string
	AllowedEmails     []string
}

func parseList(raw string) []string {
	if raw == "" {
		return nil
	}
	var values []string
	for _, v := range strings.Split(raw, ",") {
		if v = strings.TrimSpace(v); v != "" {
			values = append(values, v)
		}
	}
	return values
}

func loadAllowedEmails() []string {
	if raw := os.Getenv("OAUTH_ALLOWED_EMAILS"); raw != "" {
		return parseList(raw)
	}
	if raw := os.Getenv("OAUTH_ALLOWED_EMAIL"); raw != "" {
		return parseList(raw)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Load() *Config {
	return &Config{
		FredAPIKey:        os.Getenv("FRED_API_KEY"),
		Port:              getEnv("PORT", "4000"),
		OAuthIssuer:       getEnv("OAUTH_ISSUER", "https://accounts.google.com"),
		OAuthJwksURL:      getEnv("OAUTH_JWKS_URL", "https://www.googleapis.com/oauth2/v3/certs"),
		OAuthAuthorizeURL: getEnv("OAUTH_AUTHORIZE_URL", "https://accounts.google.com/o/oauth2/v2/auth"),
		OAuthTokenURL:     getEnv("OAUTH_TOKEN_URL", "https://oauth2.googleapis.com/token"),
		OAuthAudience:     os.Getenv("OAUTH_AUDIENCE"),
		OAuthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		OAuthCallbackURL:  os.Getenv("OAUTH_CALLBACK_URL"),
		OAuthRedirectURIs: parseList(os.Getenv("OAUTH_REDIRECT_URIS")),
		AllowedEmails:     loadAllowedEmails(),
	}
}
