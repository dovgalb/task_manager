package transport_http

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/auth/repo"
	"task-manager/internal/auth/usecases"
	"task-manager/pkg/logger/sl"
)

func DeleteHandler(log *slog.Logger, service *usecases.UserService, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	const op = "internal.handlers.rest.user.delete.DeleteHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestDelete
		_, claims, _ := jwtauth.FromContext(r.Context())
		userID := claims["user_id"].(float64)

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Ошибка декодирования запроса", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "неверный формат запроса"})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("Некорректный запрос", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "некорректные данные"})

			return
		}

		user, err := service.GetUserByID(r.Context(), userID, req.Password)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrUserNotFound):
				log.Error("Ошибка удаления пользователя", sl.Err(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, Response{Status: "error", Error: "Что-то пошло не так"})

				return
			case errors.Is(err, usecases.ErrIncorrectCredentials):
				log.Info("Не правильный логин или пароль")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, Response{Status: "error", Error: "Неверный логин или пароль"})

				return
			default:
				log.Error("неизвестная ошибка", sl.Err(err))
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, Response{Status: "error", Error: "Что-то пошло не так"})

				return
			}

		}

		err = service.DeleteUser(r.Context(), user)
		if err != nil {
			if errors.Is(err, repo.ErrUserNotFound) {
				log.Info("Такого пользователя не существует", slog.String("login", user.Login))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, Response{Status: "error", Error: "Ошибка удаления"})

			}

			log.Error("Такого пользователя не существует", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{Status: "error", Error: "Что-то пошло не так"})
			return
		}

		log.Info(fmt.Sprintf("Пользователь %s успешно удален", user.Login))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Status: "ok", Error: "Успешное удаление"})

	}
}
