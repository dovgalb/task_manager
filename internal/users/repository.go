package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"task-manager/pkg/clients/posgresql"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *User) error
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	dbClient posgresql.DBClient
	logger   *slog.Logger
}

func (r repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (login, password_hash) 
		VALUES($1, $2)
		RETURNING id
		`
	if err := r.dbClient.QueryRow(ctx, query, u.Login, u.PasswordHash).Scan(&u.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(
				fmt.Sprintf("Ошибка выполнения запроса: %s, Detail: %s, Where: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where),
			)
			return err
		}
		return err
	}

	return nil
}

func (r repository) FindAll(ctx context.Context) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindOne(ctx context.Context, id int) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(ctx context.Context, u *User) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(dbClient posgresql.DBClient, logger *slog.Logger) RepositoryInterface {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}
