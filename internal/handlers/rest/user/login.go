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

// LoginHandler эндпоинт авторизации существующего пользователя
func LoginHandler(log *slog.Logger, service *users.UserService) http.HandlerFunc {
	const op = "internal.handlers.rest.user.create.LoginHandler"
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
			render.JSON(w, r, Response{Status: "error", Error: "Что-то пошло не так"})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation error", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Ошибка валидации"})
			return
		}

		userDTO := users.UsersDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		isValid, ok := service.LoginUser(r.Context(), userDTO)
		if !ok {
			log.Error("Ошибка сервера")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "internalServerError"})
			return
		}
		if !isValid && ok {
			log.Info("неправильный пароль или логин")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, Response{Status: "error", Error: "неправильный пароль или логин"})
			return
		}
		if isValid {
			log.Info(fmt.Sprintf("Пользователь %s успешно авторизовался", userDTO.Login))
			render.Status(r, http.StatusOK)
			render.JSON(w, r, Response{Status: "ok"})
		}

	}

}
