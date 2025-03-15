package app

import (
	"github.com/go-chi/chi"
	"log/slog"
	grpcapp "task-manager/internal/app/grpc"
	httpapp "task-manager/internal/app/http"
	"task-manager/internal/config"
)

type App struct {
	GRPCSrv    *grpcapp.App
	HTTPServer *httpapp.App
}

func New(log *slog.Logger, router *chi.Mux, cnf *config.Config) *App {
	grpcApp := grpcapp.New(log, cnf)
	httpApp := httpapp.New(log, router, cnf)

	return &App{
		grpcApp,
		httpApp,
	}
}
