package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/atomic/internal/lib/logger/sl"
	"github.com/romankravchuk/atomic/internal/server/http/api/param"
	resp "github.com/romankravchuk/atomic/internal/server/http/api/response"
	"github.com/romankravchuk/atomic/internal/storage"
	"golang.org/x/exp/slog"
)

type AliasDeleter interface {
	DeleteAlias(alias string) error
}

func New(log *slog.Logger, deleter AliasDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.alias.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		name := chi.URLParam(r, param.AliasKey)
		if name == "" {
			log.Error("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := deleter.DeleteAlias(name)
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Info("alias not found", slog.String("alias", name))

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete alias", sl.Err(err))

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		log.Info("alias deleted", slog.String("alias", name))

		render.JSON(w, r, resp.OK())
	}
}
