package user

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/users"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, service *users.UserService) http.HandlerFunc {
	const op = "internal.handlers.rest.user.create.New"
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("op:", op),
			slog.String("request_id:", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Ошибка декодирования", slog.Any("err", err))
			render.JSON(w, r, "ошибка декодирования")
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("err", err))
			render.JSON(w, r, err.Error())
		}

		userDTO := users.CreateUserDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		user, err := service.CreateUser(r.Context(), userDTO)
		if err != nil {
			log.Error("Ошибка при создании пользователя", slog.Any("err", err))
			render.JSON(w, r, Response{Status: "error", Error: "ошибка при создании пользователя"})
			return
		}

		log.Info("Пользователь успешно создан", slog.Any("user", user))
		render.JSON(w, r, Response{Status: "ok"})
	}
}
