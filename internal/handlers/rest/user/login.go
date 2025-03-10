package user

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/users"
)

// LoginHandler эндпоинт авторизации существующего пользователя
func LoginHandler(log *slog.Logger, service *users.UserService, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
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

		user, err := service.AuthenticateUser(r.Context(), userDTO)
		if err != nil {
			log.Error("Ошибка аутентификации", slog.Any("err", err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, Response{Status: "error", Error: "неверный логин или пароль"})
			return
		}

		_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})

		if err != nil {
			log.Error("Ошибка генерации токена", slog.Any("err", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "Ошибка генерации токена"})
			return
		}

		log.Info("Пользователь успешно авторизован", slog.Any("user", user.Login))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Status: "ok", Token: tokenString})
	}

}
