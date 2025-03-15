package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"task-manager/internal/config"
	authgrpc "task-manager/internal/tasks_categories/transport/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, cnf *config.Config) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cnf.GRPCServer.Port,
	}

}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With("op", op)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("grpc-server запускается", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("Остановка grpc-сервера")

	a.gRPCServer.GracefulStop()
}
