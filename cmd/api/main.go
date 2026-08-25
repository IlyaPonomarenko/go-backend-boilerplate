package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"api-consume-example/internal/config"
	"api-consume-example/internal/handler"
	"api-consume-example/internal/middleware"
	"api-consume-example/internal/repository"
	"api-consume-example/internal/service"

	_ "modernc.org/sqlite"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		return
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := os.MkdirAll(filepath.Dir(settings.Database), 0750); err != nil {
		slog.Error("create data directory", "error", err)
		return
	}
	database, err := sql.Open("sqlite", settings.Database)
	if err != nil {
		slog.Error("open database", "error", err)
		return
	}
	defer database.Close()

	postRepository, err := repository.NewSQLitePostRepository(context.Background(), database)
	if err != nil {
		slog.Error("initialize database", "error", err)
		return
	}

	// main wires the application's dependencies. Keeping this composition here
	// makes the lower layers independent of HTTP and storage implementations.
	postService := service.NewPostService(postRepository)
	postHandler := handler.NewPostHandler(postService)

	mux := http.NewServeMux()
	postHandler.RegisterRoutes(mux)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("GET /healthz", handler.Health)
	apiHandler := middleware.BearerAuth(settings.APIToken)(mux)
	rootMux := http.NewServeMux()
	rootMux.Handle("/api/", apiHandler)
	rootMux.Handle("/healthz", healthMux)

	// Go 1.22+ ServeMux patterns include HTTP methods and named path values.
	// ReadHeaderTimeout limits connections that never finish their headers.
	server := &http.Server{
		Addr:              settings.Address,
		Handler:           middleware.Logging(rootMux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignals)

	go func() {
		<-shutdownSignals
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
		}
	}()

	slog.Info("API listening", "address", "http://localhost"+settings.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("server error:", err)
	}
}
