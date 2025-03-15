package httpapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"task-manager/internal/config"
	"time"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(log *slog.Logger, router *chi.Mux, cnf *config.Config) *App {

	log.Info("starting http-server at ", slog.Any("address", cnf.HTTPServer.Addr))

	server := &http.Server{
		Addr:         cnf.HTTPServer.Addr,
		Handler:      router,
		ReadTimeout:  cnf.ReadTimeout,
		WriteTimeout: cnf.WriteTimeout,
		IdleTimeout:  cnf.IdleTimeout,
	}
	return &App{
		log:        log,
		httpServer: server,
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "httpapp.Stop"
	log := a.log.With("op", op)

	// Создаем контекст с таймаутом для завершения
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Ошибка остановки HTTP-сервера", slog.Any("err", err))
	} else {
		log.Info("HTTP-сервер успешно остановлен")
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Ошибка запуска HTTP-сервера", slog.Any("err", err))
			panic("ошибка запуска http-сервера")
		}
	}
}
