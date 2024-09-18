package capture

import (
	"time"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"go.temporal.io/sdk/workflow"
)

type CapturePokemonOutput struct {
	Trainer *pokemon.Trainer
	Pokemon *pokemon.Pokemon
}

func (w *Worker) CapturePokemonWorkflow(ctx workflow.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) (*CapturePokemonOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("CapturePokemon workflow started")

	err := workflow.ExecuteActivity(ctx, w.ParalyzeActivity, pokemon).Get(ctx, pokemon)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	err = workflow.ExecuteActivity(ctx, w.AttackActivity, pokemon).Get(ctx, pokemon)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	var throwPokeballOutput ThrowPokeballOutput
	err = workflow.ExecuteActivity(ctx, w.ThrowPokeballActivity, trainer, pokemon).Get(ctx, &throwPokeballOutput)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	logger.Info("CapturePokemon workflow completed.")

	return &CapturePokemonOutput{
		Trainer: throwPokeballOutput.Trainer,
		Pokemon: throwPokeballOutput.Pokemon,
	}, nil
}
