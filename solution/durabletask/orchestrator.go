package durabletask

import (
	"github.com/microsoft/durabletask-go/task"
	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type OrchestratorInput struct {
	Trainer *pokemon.Trainer
	Pokemon *pokemon.Pokemon
}

type OrchestratorOutput struct {
	Trainer *pokemon.Trainer
	Pokemon *pokemon.Pokemon
}

func (w *Worker) CapturePokemonOrchestrator(ctx *task.OrchestrationContext) (any, error) {
	var input OrchestratorInput
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}

	var output ActivityInput
	if err := ctx.CallActivity(w.Paralyze, task.WithActivityInput(NewActivityInput(input.Trainer, input.Pokemon))).Await(&output); err != nil {
		return nil, err
	}

	if err := ctx.CallActivity(w.Attack, task.WithActivityInput(NewActivityInput(output.Trainer, output.Pokemon))).Await(&output); err != nil {
		return nil, err
	}

	if err := ctx.CallActivity(w.ThrowPokeball, task.WithActivityInput(NewActivityInput(output.Trainer, output.Pokemon))).Await(&output); err != nil {
		return nil, err
	}

	return OrchestratorOutput{
		Trainer: output.Trainer,
		Pokemon: output.Pokemon,
	}, nil
}
