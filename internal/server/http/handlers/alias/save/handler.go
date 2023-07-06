package save

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/romankravchuk/atomic/internal/data"
	"github.com/romankravchuk/atomic/internal/lib/logger/sl"
	"github.com/romankravchuk/atomic/internal/lib/random"
	resp "github.com/romankravchuk/atomic/internal/server/http/api/response"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL  string `json:"url" validate:"required,url"`
	Name string `json:"name,omitempty" validate:"omitempty,alphanum,max=50"`
}

type Response struct {
	resp.Resposne
	Alias data.Alias `json:"alias,omitempty"`
}

const aliasLength = 10

type AliasSaver interface {
	SaveAlias(alias *data.Alias) error
}

func New(log *slog.Logger, saver AliasSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.alias.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		name := req.Name
		if name == "" {
			name, err = random.NewString(aliasLength)
			if err != nil {
				log.Error("failed to generate alias", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to generate alias"))

				return
			}
		}

		alias := &data.Alias{
			URL:  req.URL,
			Name: name,
		}

		if err := saver.SaveAlias(alias); err != nil {
			log.Error("failed to save alias", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save alias"))

			return
		}

		render.JSON(w, r, Response{
			Resposne: resp.OK(),
			Alias:    *alias,
		})
	}
}
