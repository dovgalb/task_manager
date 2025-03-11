package transport_http

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/users/usecases"
)

func DeleteHandler(log *slog.Logger, service *usecases.UserService, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	const op = "internal.handlers.rest.user.delete.DeleteHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestDelete
		_, claims, _ := jwtauth.FromContext(r.Context())
		userID := claims["user_id"].(float64)

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Ошибка декодирования запроса", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "неверный формат запроса"})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("Некорректный запрос", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "некорректные данные"})

			return
		}

		user, err := service.GetUserByID(r.Context(), userID, req.Password)
		if err != nil {
			log.Error("Ошибка", slog.Any("err", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Status: "error", Error: "Неверный пароль"})

			return
		}

		if err := service.DeleteUser(r.Context(), user); err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, Response{Status: "error", Error: "Ошибка удаления"})

			return
		}

		log.Info(fmt.Sprintf("Пользователь %v успешно удален", user))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{Status: "ok", Error: "Успешное удаление"})

	}
}
