package main

import (
	"context"
	"task-manager/internal/config"
	"task-manager/pkg/clients/posgresql"
	logs "task-manager/pkg/utils"
)

func main() {
	ctx := context.Background()
	cnf := config.New()
	log := logs.SetupLogger()

	DBClient, err := posgresql.NewDBClient(ctx, cnf, log)
	if err != nil {
		log.Error("Не удалось создать клиент: error", err)
	}

}
