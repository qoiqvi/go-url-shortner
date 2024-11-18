package redirect

import (
	"log/slog"
	"net/http"

	resp "github.com/citraqs/go-url-shortner/internal/lib/api/response"
	"github.com/citraqs/go-url-shortner/internal/lib/logger/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetURL interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getUrl GetURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		url, err := getUrl.GetURL(alias)
		if err != nil {
			log.Error("Cannot get url by alias", sl.Err(err))
			render.JSON(w, r, resp.Error("Cannot get url by alias"))
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
