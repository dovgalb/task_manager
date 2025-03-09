package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"task-manager/internal/handlers/rest/user"
	"task-manager/internal/users"
)

func UsersRoutes(r *chi.Mux, log *slog.Logger, userService *users.UserService, tokenAuth *jwtauth.JWTAuth) {
	// Публичные маршруты
	r.Group(func(r chi.Router) {
		r.Post("/register", user.RegisterHandler(log, userService))
		r.Post("/login", user.LoginHandler(log, userService, tokenAuth))
	})

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))      // Ищет токен в запросе
		r.Use(jwtauth.Authenticator(tokenAuth)) // Проверяет токен

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			userID := claims["user_id"].(float64)
			render.JSON(w, r, map[string]float64{"user_id": userID})
		})
		r.Delete("/user", user.DeleteHandler(log, userService, tokenAuth))
	})
}
