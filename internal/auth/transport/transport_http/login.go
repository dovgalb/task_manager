package transport_http

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/auth/usecases"
	"task-manager/pkg/logger/sl"
)

// LoginHandler эндпоинт авторизации существующего пользователя
func LoginHandler(log *slog.Logger, service *usecases.UserService, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	const op = "internal.handlers.rest.user.create.LoginHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op:", op),
			slog.String("request_id:", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Ошибка декодирования запроса", sl.Err(err))
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

		userDTO := usecases.UsersDTO{
			Login:    req.Login,
			Password: req.Password,
		}

		user, err := service.AuthenticateUser(r.Context(), userDTO)
		if err != nil {
			if errors.Is(err, usecases.ErrIncorrectCredentials) {
				log.Info("Ошибка аутентификации")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, Response{Status: "error", Error: "неверный логин или пароль"})
				return
			}
			log.Error("Ошибка авторизации пользователя", sl.Err(err))
			return

		}

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})

		if err != nil {
			log.Error("Ошибка генерации токена", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "Ошибка генерации токена"})
			return
		}

		log.Info("Пользователь успешно авторизован", slog.Any("user", user.Login))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Status: "ok", Token: tokenString})
	}

}
