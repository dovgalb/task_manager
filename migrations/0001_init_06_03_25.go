package main

import (
	"context"
	"log/slog"
	"os"
	"task-manager/internal/config"
	"task-manager/pkg/clients/posgresql"
	"task-manager/pkg/reusable/custom_loggers"
	"time"
)

func createTasksTable(ctx context.Context, log *slog.Logger, dbClient posgresql.DBClient) error {
	const op = "migrations.0001_init_06_03_05.createTasksTable"
	stmt := `
	CREATE TABLE IF NOT EXISTS tasks(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NULL,
		is_completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	    category_id INT REFERENCES tasks_categories(id) ON DELETE SET NULL
	);
`
	_, err := dbClient.Exec(ctx, stmt)
	if err != nil {
		log.Error("Ошибка создания таблицы tasks:", err, op)
		return err
	}
	log.Info("таблица tasks успешно создалась")
	return nil
}

func createTasksCategoryTable(ctx context.Context, log *slog.Logger, dbClient posgresql.DBClient) error {
	const op = "migrations.0001_init_06_03_05.createTasksCategoryTable"

	stmt := `
	CREATE TABLE IF NOT EXISTS tasks_categories(
	    id SERIAL PRIMARY KEY,
	    title VARCHAR(255) NOT NULL UNIQUE
	);
`
	_, err := dbClient.Exec(ctx, stmt)
	if err != nil {
		log.Error("Ошибка создания таблицы tasks_categories:", err, op)
		return err
	}
	log.Info("таблица tasks_categories успешно создалась")
	return nil
}

func main() {
	log := custom_loggers.SetupLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnf := config.New()
	dbClient, err := posgresql.NewDBClient(ctx, cnf, log)

	if err != nil {
		os.Exit(1)
	}

	if err = createTasksCategoryTable(ctx, log, dbClient); err != nil {
		os.Exit(1)
	}

	if err := createTasksTable(ctx, log, dbClient); err != nil {
		os.Exit(1)
	}

	log.Info("Все таблицы успешно созданы")

}
