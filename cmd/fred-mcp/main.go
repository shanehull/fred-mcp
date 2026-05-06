package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shanehull/fred-mcp/internal/config"
	"github.com/shanehull/fred-mcp/internal/handlers"
	"github.com/shanehull/fred-mcp/internal/middleware"
	"github.com/shanehull/go-fred"
)

var (
	version = "dev"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "serve":
			runServe()
			return
		case "--version", "-version", "version":
			fmt.Println("fred-mcp", version)
			return
		}
	}

	runStdio()
}

func runStdio() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	client, err := fred.New()
	if err != nil {
		slog.Error("failed to create FRED client", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer("FRED MCP", version)
	handlers.RegisterTools(s, client)

	slog.Info("Starting FRED MCP server over stdio")
	if err := server.ServeStdio(s); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func runServe() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	client, err := fred.New(fred.WithAPIKey(cfg.FredAPIKey))
	if err != nil {
		slog.Error("failed to create FRED client", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer("FRED MCP", version)
	handlers.RegisterTools(s, client)

	var jwks *keyfunc.JWKS
	if cfg.OAuthAudience != "" {
		jwks, err = keyfunc.Get(cfg.OAuthJwksURL, keyfunc.Options{
			RefreshInterval: time.Hour,
		})
		if err != nil {
			slog.Warn("failed to fetch JWKS", "url", cfg.OAuthJwksURL, "error", err)
		}
	}

	sse := server.NewSSEServer(s, server.WithBaseURL(cfg.PublicHost))
	streamable := server.NewStreamableHTTPServer(s)

	auth := middleware.Auth(cfg, jwks)

	mux := http.NewServeMux()
	registerHTTPHandlers(mux, cfg, sse, streamable, auth)

	handler := loggingMiddleware(enableCORS(mux))

	slog.Info("FRED MCP server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func registerHTTPHandlers(mux *http.ServeMux, cfg *config.Config, sse *server.SSEServer, streamable *server.StreamableHTTPServer, auth func(http.Handler) http.Handler) {
	mux.Handle("/sse", auth(sse.SSEHandler()))
	mux.Handle("/message", auth(sse.MessageHandler()))
	mux.Handle("/mcp", auth(streamable))

	mux.Handle("/.well-known/oauth-protected-resource/", handlers.HandleDiscovery(cfg))
	mux.HandleFunc("/.well-known/oauth-protected-resource", handlers.HandleDiscovery(cfg))
	mux.HandleFunc("/.well-known/oauth-authorization-server", handlers.HandleAuthServerDiscovery(cfg))
	mux.HandleFunc("/.well-known/mcp", handlers.HandleDiscovery(cfg))

	mux.HandleFunc("/register", handlers.HandleRegistration(cfg))
	mux.HandleFunc("/authorize", handlers.HandleAuthorizeProxy(cfg))
	mux.HandleFunc("/token", handlers.HandleTokenProxy(cfg))
	mux.HandleFunc("/config", handlers.HandleConfig(cfg))
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
			"duration", time.Since(start),
		)
	})
}
