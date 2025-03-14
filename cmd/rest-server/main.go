package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-manager/internal/auth/repo"
	"task-manager/internal/auth/transport/transport_http"
	"task-manager/internal/auth/usecases"
	"task-manager/internal/config"
	"task-manager/pkg/clients/kafka"
	"task-manager/pkg/clients/posgresql"
	"task-manager/pkg/logger/handlers/slogpretty"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cnf := config.New()
	log := SetupLogger()

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

	log.Info("starting http-server at ", slog.Any("address", cnf.HTTPServer.Addr))
	server := &http.Server{
		Addr:         cnf.Addr,
		Handler:      router,
		ReadTimeout:  cnf.ReadTimeout,
		WriteTimeout: cnf.WriteTimeout,
		IdleTimeout:  cnf.IdleTimeout,
	}

	// асинхронно запускаем http сервер
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Ошибка запуска HTTP-сервера", slog.Any("err", err))
			}
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("Получен сигнал завершения приложения", slog.String("signal", sign.String()))

	// Создаем контекст с таймаутом для завершения
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Остановка HTTP-сервера
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Ошибка остановки HTTP-сервера", slog.Any("err", err))
	} else {
		log.Info("HTTP-сервер успешно остановлен")
	}

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
func SetupLogger() *slog.Logger {
	log := setupPrettySlog()
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
