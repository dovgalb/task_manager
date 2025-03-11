package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"
	"os"
	"task-manager/internal/config"
	"task-manager/internal/users/repo"
	"task-manager/internal/users/transport/transport_http"
	"task-manager/internal/users/usecase"
	"task-manager/pkg/clients/posgresql"
	"task-manager/pkg/kafka"
	logs "task-manager/pkg/utils"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cnf := config.New()
	log := logs.SetupLogger()

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиент: error", err)
	}

	producer, err := kafka.NewKafkaProducer(log, cnf.Brokers, cnf.Topic)
	if err != nil {
		log.Error("Ошибка продюсера", slog.Any("err", err))
	}
	defer func(producer *kafka.Producer) {
		err := producer.Close()
		if err != nil {
			log.Error("Ошибка закрытия продюсера", slog.Any("err", err))
		}
	}(producer)

	userRepository := repo.NewRepository(DBClient, log)
	userService := usecase.NewUserService(log, userRepository, producer)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	transport_http.UsersRoutes(router, log, userService, tokenAuth)

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
