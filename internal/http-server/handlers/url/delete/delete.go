package delete

import (
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "github.com/go-chi/render"
  "log/slog"
  "net/http"
  resp "url-shortener/internal/lib/api/response"
  "url-shortener/internal/lib/logger/sl"
  "url-shortener/internal/storage/sqlite"
)

type Response struct {
  resp.Response
  Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, storage *sqlite.Storage) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.url.delete.New"
    log = log.With(
      slog.String("op", op),
      slog.String("request_id", middleware.GetReqID(r.Context())),
    )
    alias := chi.URLParam(r, "alias")
    if err := storage.DeleteURL(alias); err != nil {
      log.Error("failed to delete url", sl.Err(err))
      render.JSON(w, r, resp.Error("failed to delete url"))
      return
    }
    responseOK(w, r, alias)
  }
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
  render.JSON(w, r, Response{
    Response: resp.OK(),
    Alias:    alias,
  })
}
