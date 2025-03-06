package main

import (
	"context"
	"task-manager/internal/config"
	"task-manager/pkg/clients/posgresql"
	logs "task-manager/pkg/reusable/custom_loggers"
)

func main() {
	ctx := context.Background()
	cnf := config.New()
	log := logs.SetupLogger()

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиент: error", err)
	}

	stmt := `
CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
`
	ct, err := DBClient.Exec(ctx, stmt)

}
