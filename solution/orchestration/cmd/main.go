package main

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/solution/orchestration"
)

func main() {
	combatService := pokemon.NewCombatService()
	statusService := pokemon.NewStatusService()
	pokeballService := pokemon.NewPokeballService()

	orchestrator := orchestration.NewOrchestrator()

	orchestrator.Register(orchestration.NewCapturePokemonWorker(
		statusService,
		combatService,
		pokeballService,
	))

	workflowExecution := orchestration.NewCapturePokemonWorkflowExecution(pokemon.Sacha(), pokemon.Mewtwo())

	err := orchestrator.RunWorkflow(workflowExecution)
	if err != nil {
		slog.Error("fail to run workflow", slog.Any("error", err))
	}
}
