package tasks_category

import (
	"context"
	"log/slog"
	"task-manager/pkg/clients/posgresql"
)

type Repository interface {
	Create(ctx context.Context, task TaskCategory) (string, error)
	FindAll(ctx context.Context) ([]TaskCategory, error)
	FindOne(ctx context.Context, id int) (TaskCategory, error)
	Update(ctx context.Context, task TaskCategory) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	dbClient posgresql.DBClient
	logger   *slog.Logger
}

func (r repository) Create(ctx context.Context, task TaskCategory) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindAll(ctx context.Context) ([]TaskCategory, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindOne(ctx context.Context, id int) (TaskCategory, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ctx context.Context, task TaskCategory) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(dbClient posgresql.DBClient, logger *slog.Logger) Repository {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}
