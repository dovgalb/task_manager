package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/handlers/rest/user"
	"task-manager/internal/users"
	"task-manager/pkg/clients/posgresql"
	logs "task-manager/pkg/utils"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cnf := config.New()
	log := logs.SetupLogger()

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиент: error", err)
	}

	userRepository := users.NewRepository(DBClient, log)
	userService := users.NewUserService(userRepository, log)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Защищенные маршруты
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			userID := claims["user_id"].(float64)
			render.JSON(w, r, map[string]float64{"user_id": userID})
		})
	})

	// Публичные маршруты
	router.Group(func(r chi.Router) {
		r.Post("/register", user.RegisterHandler(log, userService))
		r.Post("/login", user.LoginHandler(log, userService, tokenAuth))
	})

	log.Info("starting http-server at ", slog.Any("address", cnf.HTTPServer.Addr))
	server := &http.Server{
		Addr:         cnf.Addr,
		Handler:      router,
		ReadTimeout:  cnf.ReadTimeout,
		WriteTimeout: cnf.WriteTimeout,
		IdleTimeout:  cnf.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
		os.Exit(1)
	}

	log.Error("server stopped")
}
