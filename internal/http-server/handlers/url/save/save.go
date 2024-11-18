package save

import (
	"log/slog"
	"net/http"

	resp "github.com/citraqs/go-url-shortner/internal/lib/api/response"
	"github.com/citraqs/go-url-shortner/internal/lib/logger/sl"
	"github.com/citraqs/go-url-shortner/internal/lib/random"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.BaseResponse
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save"
		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))

			return
		}

		log.Info("Request Body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.GenerateRandomString(aliasLength)
		}

		_, err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error("Failed to insert data in DB", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to insert data"))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
