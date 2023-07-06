package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/atomic/internal/config"
	"github.com/romankravchuk/atomic/internal/lib/logger/sl"
	"github.com/romankravchuk/atomic/internal/server/http/handlers/alias/delete"
	"github.com/romankravchuk/atomic/internal/server/http/handlers/alias/redirect"
	"github.com/romankravchuk/atomic/internal/server/http/handlers/alias/save"
	mwLogger "github.com/romankravchuk/atomic/internal/server/middleware/logger"
	"github.com/romankravchuk/atomic/internal/storage/postgresql"
	"github.com/romankravchuk/atomic/internal/storage/postgresql/alias"
	"golang.org/x/exp/slog"
)

const (
	localEnv = "local"
	prodEnv  = "prod"
	devEnv   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg)

	log.Info("starting up...")
	log.Debug("debug mode activated")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)

	db, err := postgresql.New(dsn)
	if err != nil {
		log.Error("failed to open database connection", sl.Err(err))
		os.Exit(1)
	}

	storage, err := alias.New(db)
	if err != nil {
		log.Error("failed to open storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api", func(r chi.Router) {
		r.Route("/alias", func(r chi.Router) {
			r.Use(middleware.BasicAuth("atomic", map[string]string{
				cfg.HttpServer.Username: cfg.HttpServer.Password,
			}))
			r.Post("/", save.New(log, storage))
			r.Delete("/{alias}", delete.New(log, storage))
		})
		r.Get("/{alias}", redirect.New(log, storage))
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.HttpServer.Addr,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}

	log.Error("server stopped")
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case localEnv:
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case devEnv:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prodEnv:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
