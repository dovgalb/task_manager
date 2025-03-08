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
	UserID int    `json:"user_id,omitempty"`
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
			log.Error("Ошибка декодирования запроса", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Неверный формат запроса"})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Некорректные данные"})
			return
		}

		userDTO := users.CreateUserDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		user, err := service.CreateUser(r.Context(), userDTO)
		if err != nil {
			log.Error("Ошибка при создании пользователя", slog.Any("err", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "ошибка при создании пользователя"})
			return
		}
		log.Info("Пользователь успешно создан", slog.Any("user", user))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{Status: "ok", UserID: user.ID})
	}
}
