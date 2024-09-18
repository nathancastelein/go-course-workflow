package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/microsoft/durabletask-go/api"
	"github.com/microsoft/durabletask-go/backend"
	"github.com/microsoft/durabletask-go/backend/sqlite"
	"github.com/microsoft/durabletask-go/task"
	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/durabletask"
)

func main() {
	worker := durabletask.NewWorker(
		pokemon.NewStatusService(),
		pokemon.NewCombatService(),
		pokemon.NewPokeballService(),
	)

	logger := newLogger()

	ctx := context.Background()

	r := task.NewTaskRegistry()
	r.AddOrchestrator(worker.CapturePokemonOrchestrator)
	r.AddActivity(worker.Paralyze)
	r.AddActivity(worker.Attack)
	r.AddActivity(worker.ThrowPokeball)

	// Create an executor
	executor := task.NewTaskExecutor(r)

	// Create a new backend
	// Use the in-memory sqlite provider by specifying ""
	be := sqlite.NewSqliteBackend(sqlite.NewSqliteOptions(""), logger)
	orchestrationWorker := backend.NewOrchestrationWorker(be, executor, logger)
	activityWorker := backend.NewActivityTaskWorker(be, executor, logger)
	taskHubWorker := backend.NewTaskHubWorker(be, orchestrationWorker, activityWorker, logger)

	// Start the worker
	err := taskHubWorker.Start(ctx)
	if err != nil {
		slog.Error("fail to start worker", slog.Any("error", err))
		return
	}

	// Get the client to the backend
	taskHubClient := backend.NewTaskHubClient(be)

	ctx = context.Background()
	defer taskHubWorker.Shutdown(ctx)

	// Start a new orchestration
	id, err := taskHubClient.ScheduleNewOrchestration(ctx, worker.CapturePokemonOrchestrator, api.WithInput(durabletask.OrchestratorInput{
		Trainer: pokemon.Sacha(),
		Pokemon: pokemon.Mewtwo(),
	}))
	if err != nil {
		slog.Error("failed to schedule new orchestration", slog.Any("error", err))
		return
	}

	// Wait for the orchestration to complete
	metadata, err := taskHubClient.WaitForOrchestrationCompletion(ctx, id)
	if err != nil {
		slog.Error("failed to wait for orchestration to complete", slog.Any("error", err))
		return
	}

	var output durabletask.OrchestratorOutput
	if err := json.Unmarshal([]byte(metadata.SerializedOutput), &output); err != nil {
		slog.Error("failed to decode result from JSON", slog.Any("error", err))
		return
	}

	slog.Info("orchestration completed", slog.Any("trainer", output.Trainer), slog.Any("pokemon", output.Pokemon))
}
