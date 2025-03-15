package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"task-manager/internal/app"
	"task-manager/internal/auth/repo"
	"task-manager/internal/auth/transport/transport_http"
	"task-manager/internal/auth/usecases"
	"task-manager/internal/config"
	"task-manager/pkg/clients/kafka"
	"task-manager/pkg/clients/posgresql"
	"task-manager/pkg/logger/handlers/slogpretty"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cnf := config.New()
	log := SetupLogger(cnf.Env)

	log.Info("Приложение запущено", slog.String("Окружение", cnf.Env))

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиента базы данных: error", err)
	}

	producer, err := kafka.NewKafkaProducer(log, cnf.Brokers, cnf.Topic)
	if err != nil {
		log.Error("Ошибка создания Kafka продюсера", slog.Any("err", err))
	}

	userRepository := repo.NewRepository(DBClient)
	userService := usecases.NewUserService(log, userRepository, producer)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	transport_http.UsersRoutes(router, log, userService, tokenAuth)

	application := app.New(log, router, cnf)
	go application.GRPCSrv.MustRun()
	go application.HTTPServer.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("Получен сигнал завершения приложения", slog.String("signal", sign.String()))

	// Остановка HTTP-сервера
	application.HTTPServer.Stop()
	// остановка GRPC-сервера
	application.GRPCSrv.Stop()

	// Закрытие Kafka producer
	if err := producer.Close(); err != nil {
		log.Error("Ошибка остановки Kafka-продюсера", slog.Any("err", err))
	} else {
		log.Info("Kafka-продюсер успешно остановлен")
	}

	// Закрытие клиента Базы данных
	DBClient.Close()
	log.Info("Клиент базы данных закрыт")

	log.Info("Программа завершена")

	// TODO GRPC ручки для tasks и tasks_categories
	// TODO реализовать нормальные миграции
	// TODO подтверждение сообщения после прочтения
	// TODO написать тесты для ручек с моками(mockery)
	// TODO написать функциональные тесты
	// TODO переработать получение токена
}

// SetupLogger Устанавливает логгер
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = setupPrettySlog()
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
