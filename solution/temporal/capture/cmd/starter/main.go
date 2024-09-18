package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/temporal/capture"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

func main() {
	c, err := client.Dial(client.Options{
		Logger: log.NewStructuredLogger(slog.Default()),
	})
	if err != nil {
		slog.Error("unable to create client", slog.Any("error", err))
		return
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "capture-pokemon",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, "CapturePokemonWorkflow", pokemon.Sacha(), pokemon.Mewtwo())
	if err != nil {
		slog.Error("unable to execute workflow", slog.Any("error", err))
		return
	}

	slog.Info("started workflow", slog.String("workflow_id", we.GetID()), slog.String("run_id", we.GetRunID()))

	var result capture.CapturePokemonOutput
	err = we.Get(context.Background(), &result)
	if err != nil {
		slog.Error("unable to get workflow result", slog.Any("error", err))
		return
	}
	slog.Info("workflow result", slog.Any("trainer", result.Trainer), slog.Any("pokemon", result.Pokemon))
}
