package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"task-manager/pkg/clients/posgresql"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *User) error
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, login string) (*User, error)
	FindOneByID(ctx context.Context, id int) (*User, error)
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

func (r repository) FindOne(ctx context.Context, login string) (*User, error) {
	query := `
	SELECT id, login, password_hash, created_at, updated_at
	FROM users 
	WHERE login = $1
`

	var user User
	err := r.dbClient.QueryRow(ctx, query, login).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("Пользователь не найден", slog.String("login", login))
			return &User{}, fmt.Errorf("Пользователь %s не найден, %w", login, err)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(
				fmt.Sprintf("Ошибка выполнения запроса: %s, Detail: %s, Where: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where),
			)
		}
		r.logger.Error("Ошибка", slog.Any("err", err))
		return nil, err
	}

	return &user, nil

}

func (r repository) FindOneByID(ctx context.Context, id int) (*User, error) {
	query := `
	SELECT id, login, password_hash, created_at, updated_at
	FROM users 
	WHERE id = $1
`
	var user User
	err := r.dbClient.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("Пользователь не найден", slog.Int("id", id))
			return &User{}, fmt.Errorf("пользователь %d не найден, %w", id, err)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(
				fmt.Sprintf("Ошибка выполнения запроса: %s, Detail: %s, Where: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where),
			)
		}
		r.logger.Error("Ошибка", slog.Any("err", err))
		return nil, err
	}

	return &user, nil

}

func (r repository) Update(ctx context.Context, u *User) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE 
	FROM users 
	WHERE id = $1
`
	pgTag, err := r.dbClient.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("ошибка удаления", slog.Any("err", err), slog.Any("operation tag", pgTag))
		return err
	}
	return nil
}

func NewRepository(dbClient posgresql.DBClient, logger *slog.Logger) RepositoryInterface {
	return &repository{
		dbClient: dbClient,
		logger:   logger,
	}
}
