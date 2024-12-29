package getAll

import (
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
  URLs []map[string]string `json:"urls,omitempty"`
}

func New(log *slog.Logger, storage *sqlite.Storage) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.url.getAll.New"
    log = log.With(
      slog.String("op", op),
      slog.String("request_id", middleware.GetReqID(r.Context())),
    )

    urls, err := storage.GetURLs()
    if err != nil {
      log.Error("failed to get urls", sl.Err(err))
      render.JSON(w, r, resp.Error("failed to get urls"))
      return
    }
    responseOK(w, r, urls)
  }
}

func responseOK(w http.ResponseWriter, r *http.Request, urls []map[string]string) {
  render.JSON(w, r, Response{
    Response: resp.OK(),
    URLs:     urls,
  })
}
