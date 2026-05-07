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
	PublicHost        string
	OAuthAudience     string
	OAuthClientSecret string
	AllowedEmails     []string
}

func parseEmails(raw string) []string {
	if raw == "" {
		return nil
	}
	var emails []string
	for _, e := range strings.Split(raw, ",") {
		e = strings.TrimSpace(e)
		if e != "" {
			emails = append(emails, e)
		}
	}
	return emails
}

func loadAllowedEmails() []string {
	if raw := os.Getenv("OAUTH_ALLOWED_EMAILS"); raw != "" {
		return parseEmails(raw)
	}
	if raw := os.Getenv("OAUTH_ALLOWED_EMAIL"); raw != "" {
		return parseEmails(raw)
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
		PublicHost:        os.Getenv("PUBLIC_HOST"),
		OAuthAudience:     os.Getenv("OAUTH_AUDIENCE"),
		OAuthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		AllowedEmails:     loadAllowedEmails(),
	}
}
