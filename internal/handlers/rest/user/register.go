package user

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/users"
)

// RegisterHandler эндпоинт регистрации нового пользователя
func RegisterHandler(log *slog.Logger, service *users.UserService) http.HandlerFunc {
	const op = "internal.handlers.rest.user.create.RegisterHandler"
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

		userDTO := users.UsersDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		user, err := service.RegisterUser(r.Context(), userDTO)
		if err != nil {
			log.Error("Ошибка при создании пользователя", slog.Any("err", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "ошибка при создании пользователя"})
			return
		}
		log.Info(fmt.Sprintf("Пользователь %s успешно создан", user.Login))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{Status: "ok", UserID: user.ID})
	}
}
