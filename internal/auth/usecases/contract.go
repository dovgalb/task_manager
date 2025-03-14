// go generate
package usecases

import (
	"context"
	"task-manager/internal/auth/repo"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *repo.User) error
	FindOne(ctx context.Context, login string) (*repo.User, error)
	FindOneByID(ctx context.Context, id int) (*repo.User, error)
	Delete(ctx context.Context, id int) error
}

type Producer interface {
	SendMessage(key, value string) error
}
