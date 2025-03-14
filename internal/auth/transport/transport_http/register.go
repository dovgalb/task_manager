package transport_http

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/auth/repo"
	"task-manager/internal/auth/usecases"
	"task-manager/pkg/logger/sl"
)

// RegisterHandler эндпоинт регистрации нового пользователя
func RegisterHandler(log *slog.Logger, service *usecases.UserService) http.HandlerFunc {
	const op = "internal.handlers.rest.user.create.RegisterHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("op", op),
			slog.String("request_id:", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Ошибка декодирования запроса", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Неверный формат запроса"})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Некорректные данные"})
			return
		}

		userDTO := usecases.UsersDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		user, err := service.RegisterUser(r.Context(), userDTO)
		if err != nil {
			if errors.Is(err, repo.ErrUserExists) {
				log.Info("Пользователь уже существует", slog.String("login", req.Login))
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, Response{Status: "error", Error: "Пользователь с таким логином уже существует"})
				return
			}

			log.Error("Ошибка при создании пользователя", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "ошибка при создании пользователя"})
			return
		}
		log = log.With(slog.String("login", userDTO.Login))

		log.Info("Пользователь успешно создан")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{Status: "ok", UserID: user.ID})
	}
}
