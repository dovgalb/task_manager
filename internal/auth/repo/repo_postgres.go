package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"task-manager/pkg/clients/posgresql"
)

var (
	ErrUserExists   = errors.New("пользователь с таким логином уже существует")
	ErrUserNotFound = errors.New("пользователь с таким логином не найден")
)

// wrapError — вспомогательная функция для обработки ошибок
func wrapError(op string, err error) error {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case "23505": // Unique constraint violation
			return fmt.Errorf("%s: %w", op, ErrUserExists)
		default:
			return fmt.Errorf("%s: %s: %w", op, pgErr.Code, err)
		}
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}

type Repository struct {
	dbClient posgresql.DBClient
}

func (r Repository) Create(ctx context.Context, u *User) error {
	const op = "auth.repo.Create"
	stmt := `
		INSERT INTO users (login, password_hash) 
		VALUES($1, $2)
		RETURNING id
		`
	err := r.dbClient.QueryRow(ctx, stmt, u.Login, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		return wrapError(op, err)
	}

	return nil
}

func (r Repository) FindAll(ctx context.Context) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) FindOne(ctx context.Context, login string) (*User, error) {
	const op = "auth.repo.FindOne"

	stmt := `
	SELECT id, login, password_hash, created_at, updated_at
	FROM users 
	WHERE login = $1
`
	var user User

	err := r.dbClient.QueryRow(ctx, stmt, login).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(op, err)
	}

	return &user, nil

}

func (r Repository) FindOneByID(ctx context.Context, id int) (*User, error) {
	const op = "auth.repo.FindOneByID"

	stmt := `
	SELECT id, login, password_hash, created_at, updated_at
	FROM users 
	WHERE id = $1
`
	var user User
	err := r.dbClient.QueryRow(ctx, stmt, id).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, wrapError(op, err)
	}

	return &user, nil
}

func (r Repository) Update(ctx context.Context, u *User) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Delete(ctx context.Context, id int) error {
	const op = "auth.repo.Delete"

	query := `
	DELETE 
	FROM users 
	WHERE id = $1
`
	pgTag, err := r.dbClient.Exec(ctx, query, id)
	if err != nil {
		return wrapError(op, err)
	}
	if pgTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return nil
}

func NewRepository(dbClient posgresql.DBClient) *Repository {
	return &Repository{
		dbClient: dbClient,
	}
}
