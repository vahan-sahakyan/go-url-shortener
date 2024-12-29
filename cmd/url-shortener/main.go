package main

import (
  "fmt"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "log/slog"
  "os"
  "url-shortener/internal/config"
  mwLogger "url-shortener/internal/http-server/middleware/logger"
  "url-shortener/internal/lib/logger/handlers/slogpretty"
  "url-shortener/internal/lib/logger/sl"
  "url-shortener/internal/storage/sqlite"
)

const (
  envLocal = "local"
  envDev   = "dev"
  envProd  = "prod"
)

func main() {
  cfg := config.MustLoad()

  fmt.Println(cfg)
  log := setupLogger(cfg.Env)

  log.Info("starting url-shortener", slog.String("env", cfg.Env), slog.String("version", "v0.1.3"))
  log.Debug("debug messages are enabled")
  log.Error("error messages are enabled")

  storage, err := sqlite.New(cfg.StoragePath)
  _ = storage
  if err != nil {
    log.Error("failed to init storage", sl.Err(err))
    os.Exit(1)
  }
  router := chi.NewRouter()

  router.Use(middleware.RequestID)
  router.Use(middleware.Logger)
  router.Use(mwLogger.New(log))
  router.Use(middleware.Recoverer)
  router.Use(middleware.URLFormat)

  // TODO: run server
}

func setupLogger(env string) *slog.Logger {
  var log *slog.Logger
  switch env {
  case envLocal:
    log = setupPrettySlog()
  case envDev:
    log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
  case envProd:
    log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
  }
  return log
}

func setupPrettySlog() *slog.Logger {
  opts := slogpretty.PrettyHandlerOptions{
    SlogOpts: &slog.HandlerOptions{
      Level: slog.LevelDebug,
    },
  }

  handler := opts.NewPrettyHandler(os.Stdout)

  return slog.New(handler)
}
