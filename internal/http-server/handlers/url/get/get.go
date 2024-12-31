package get

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
  Alias    string `json:"alias,omitempty"`
  URL      string `json:"url,omitempty"`
  Redirect string `json:"redirect,omitempty"`
}

func New(log *slog.Logger, storage *sqlite.Storage) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    const op = "handlers.url.get.New"
    log = log.With(
      slog.String("op", op),
      slog.String("request_id", middleware.GetReqID(r.Context())),
    )

    alias := chi.URLParam(r, "alias")

    url, err := storage.GetURL(alias)
    if err != nil {
      log.Error("failed to get url", sl.Err(err))
      render.JSON(w, r, resp.Error("failed to get url"))
      return
    }
    responseOK(w, r, url, alias)
  }
}

func responseOK(w http.ResponseWriter, r *http.Request, url string, alias string) {
  render.JSON(w, r, Response{
    Response: resp.OK(),
    URL:      url,
    Redirect: "http://localhost:8082/" + alias,
    Alias:    alias,
  })
}
