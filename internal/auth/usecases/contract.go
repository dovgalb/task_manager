// go generate
package usecases

import (
	"context"
	"task-manager/internal/auth/repo"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *repo.User) error
	FindAll(ctx context.Context) ([]repo.User, error)
	FindOne(ctx context.Context, login string) (*repo.User, error)
	FindOneByID(ctx context.Context, id int) (*repo.User, error)
	Update(ctx context.Context, u *repo.User) error
	Delete(ctx context.Context, id int) error
}

type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	With(args ...any) *Logger
}

type Producer interface {
	SendMessage(key, value string) error
}
