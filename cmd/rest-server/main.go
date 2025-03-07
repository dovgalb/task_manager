package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/handlers/rest/user"
	"task-manager/internal/users"
	"task-manager/pkg/clients/posgresql"
	logs "task-manager/pkg/utils"
)

func main() {
	ctx := context.Background()
	cnf := config.New()
	log := logs.SetupLogger()

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиент: error", err)
	}

	userRepository := users.NewRepository(DBClient, log)
	userService := users.NewUserService(userRepository)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", user.New(log, userService))

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
