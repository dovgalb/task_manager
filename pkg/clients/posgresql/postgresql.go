package posgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"task-manager/internal/config"
	"time"
)

type DBClient interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewDBClient(ctx context.Context, cnf *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s", cnf.DbUser, cnf.DbPassword, cnf.DbHost, cnf.DbPort, cnf.DbName,
	)

	for range cnf.DbMaxAttempts {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Error("failed to connect to postgresql")
			return nil, err
		}
		return pool, nil
	}

	return pool, nil
}
