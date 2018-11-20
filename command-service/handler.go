package main

import (
	"github.com/s-take/goecho-nats-postgre-sample/db"
	"github.com/s-take/goecho-nats-postgre-sample/event"
	"github.com/s-take/goecho-nats-postgre-sample/schema"
	"go.uber.org/zap"
)

func onTaskPublished(t event.TaskPublishedMessage) {
	logger, _ := zap.NewProduction()
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	sugar := logger.Sugar()

	task := schema.Task{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	}
	sugar.Info("Insert task: " + task.Name)
	if err := db.InsertTask(task); err != nil {
		sugar.Errorf("Internal error: %v", err)
	}
}
